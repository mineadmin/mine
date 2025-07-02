package downloader

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mineadmin/mine/internal/prompt"
)

const (
	baseURL = "https://github.com/mineadmin/mineadmin/archive/refs/tags"
)

type Downloader struct {
	Language string
	Version  string
	Platform string
}

func NewDownloader(language, version, platform string) *Downloader {
	return &Downloader{
		Language: language,
		Version:  version,
		Platform: platform,
	}
}

func (d *Downloader) Download(projectName string) error {
	// For PHP projects, download the source zip from GitHub releases
	if d.Language == "php" {
		url := fmt.Sprintf("%s/%s.zip", baseURL, d.Version)

		prompt.Info("Creating project directory...")
		spinner := prompt.StartSpinner("Setting up project structure")
		if err := os.MkdirAll(projectName, 0755); err != nil {
			spinner.Stop()
			return fmt.Errorf("failed to create project directory: %v", err)
		}
		spinner.Stop()

		// Download the file
		prompt.Info("Downloading project files...")
		spinner = prompt.StartSpinner("Fetching MineAdmin source code")
		resp, err := http.Get(url)
		if err != nil {
			spinner.Stop()
			return fmt.Errorf("failed to download: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			spinner.Stop()
			return fmt.Errorf("download failed with status: %s", resp.Status)
		}

		// Create the output file
		outputPath := filepath.Join(projectName, fmt.Sprintf("%s.zip", d.Version))
		out, err := os.Create(outputPath)
		if err != nil {
			spinner.Stop()
			return fmt.Errorf("failed to create output file: %v", err)
		}
		defer out.Close()

		// Write the body to file
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			spinner.Stop()
			return fmt.Errorf("failed to write file: %v", err)
		}
		spinner.Stop()
		prompt.Success("Download completed")

		// Unzip the file
		prompt.Info("Extracting project files...")
		spinner = prompt.StartSpinner("Unpacking MineAdmin source code")
		if err := unzip(outputPath, projectName); err != nil {
			spinner.Stop()
			return fmt.Errorf("failed to unzip: %v", err)
		}
		spinner.Stop()
		prompt.Success("Extraction completed")

		return nil
	}

	// Original download logic for other languages
	url := fmt.Sprintf("%s/%s/mineadmin-%s-%s.zip", baseURL, d.Version, d.Language, d.Platform)

	prompt.Info("Creating project directory...")
	spinner := prompt.StartSpinner("Setting up project structure")
	if err := os.MkdirAll(projectName, 0755); err != nil {
		spinner.Stop()
		return fmt.Errorf("failed to create project directory: %v", err)
	}
	spinner.Stop()

	// Download the file
	prompt.Info("Downloading project files...")
	spinner = prompt.StartSpinner("Fetching MineAdmin source code")
	resp, err := http.Get(url)
	if err != nil {
		spinner.Stop()
		return fmt.Errorf("failed to download: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		spinner.Stop()
		return fmt.Errorf("download failed with status: %s", resp.Status)
	}

	// Create the output file
	outputPath := filepath.Join(projectName, fmt.Sprintf("mineadmin-%s-%s.zip", d.Language, d.Platform))
	out, err := os.Create(outputPath)
	if err != nil {
		spinner.Stop()
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		spinner.Stop()
		return fmt.Errorf("failed to write file: %v", err)
	}
	spinner.Stop()
	prompt.Success("Download completed")

	return nil
}

func (d *Downloader) ListVersions() ([]string, error) {
	if d.Language == "php" {
		resp, err := http.Get("https://api.github.com/repos/mineadmin/mineadmin/releases")
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

	// For other languages, return mock data
	return []string{"v1.0.0", "v1.0.1", "v1.1.0"}, nil
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		// Skip directories
		if f.FileInfo().IsDir() {
			continue
		}

		// Remove the top-level directory from the path
		parts := strings.Split(f.Name, string(filepath.Separator))
		if len(parts) < 2 {
			continue // Skip files in root directory
		}
		relPath := filepath.Join(parts[1:]...)

		// Open the file inside the zip
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		// Create the output file
		path := filepath.Join(dest, relPath)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		out, err := os.Create(path)
		if err != nil {
			return err
		}
		defer out.Close()

		// Copy the contents
		_, err = io.Copy(out, rc)
		if err != nil {
			return err
		}
	}

	return nil
}
