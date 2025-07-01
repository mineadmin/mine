package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/mineadmin/mine/internal/downloader"
	"github.com/spf13/cobra"
)

// VersionQuerier 定义版本查询接口
type VersionQuerier interface {
	ListVersions() ([]string, error)
}

// DownloaderQuerier 实现原下载器的版本查询
type DownloaderQuerier struct {
	language string
}

func (q *DownloaderQuerier) ListVersions() ([]string, error) {
	dl := downloader.NewDownloader(q.language, "", "")
	return dl.ListVersions()
}

// GitHubQuerier 实现GitHub API的版本查询
type GitHubQuerier struct {
	repo string
}

func (q *GitHubQuerier) ListVersions() ([]string, error) {
	resp, err := http.Get("https://api.github.com/repos/" + q.repo + "/releases")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var releases []struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, err
	}

	var versions []string
	for _, r := range releases {
		versions = append(versions, r.TagName)
	}
	return versions, nil
}

// NewVersionQuerier 工厂方法创建版本查询器
func NewVersionQuerier(language string) VersionQuerier {
	switch language {
	case "php":
		return &GitHubQuerier{repo: "mineadmin/mineadmin"}
	default:
		return &DownloaderQuerier{language: language}
	}
}

// NewSelectVersionsCmd creates and returns the select-versions command
func NewSelectVersionsCmd() *cobra.Command {
	var language string

	cmd := &cobra.Command{
		Use:   "select-versions",
		Short: "List available versions of MineAdmin",
		Long: `List all available versions of MineAdmin for specified language.
Example:
  mine select-versions --language=php`,
		Run: func(cmd *cobra.Command, args []string) {
			querier := NewVersionQuerier(language)
			versions, err := querier.ListVersions()
			if err != nil {
				log.Fatalf("Failed to list versions: %v", err)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "VERSION\tLANGUAGE\tSTATUS")
			for _, v := range versions {
				fmt.Fprintf(w, "%s\t%s\t%s\n", v, language, "available")
			}
			w.Flush()
		},
	}

	cmd.Flags().StringVarP(&language, "language", "l", "", "Programming language (required)")
	cmd.MarkFlagRequired("language")

	return cmd
}
