package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "mine",
		Short: "MineAdmin CLI tool",
		Long: `A command line tool for downloading and managing MineAdmin projects.
Complete documentation is available at https://github.com/mineadmin/mine`,
	}

	// Add all subcommands
	rootCmd.AddCommand(NewCreateCmd())
	rootCmd.AddCommand(NewSelectVersionsCmd())

	return rootCmd
}

var rootCmd = NewRootCmd()

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
