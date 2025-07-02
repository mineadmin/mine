package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	var binPhp string

	rootCmd := &cobra.Command{
		Use:   "mine",
		Short: "MineAdmin CLI tool",
		Long: `
╭───────────────────────────────────────────────╮
│                                               │
│   ███╗   ███╗██╗███╗   ██╗███████╗           │
│   ████╗ ████║██║████╗  ██║██╔════╝           │
│   ██╔████╔██║██║██╔██╗ ██║█████╗             │
│   ██║╚██╔╝██║██║██║╚██╗██║██╔══╝             │
│   ██║ ╚═╝ ██║██║██║ ╚████║███████╗           │
│   ╚═╝     ╚═╝╚═╝╚═╝  ╚═══╝╚══════╝           │
│                                               │
│   MineAdmin CLI - 项目创建和管理工具         │
│                                               │
╰───────────────────────────────────────────────╯

A powerful command line tool for downloading and managing MineAdmin projects.

🔹 Commands:
  - create: Create a new MineAdmin project
  - select-versions: List available MineAdmin versions

🔹 Examples:
  mine create my-project
  mine select-versions

Complete documentation is available at https://github.com/mineadmin/mine`,
	}

	// Add global flags
	rootCmd.PersistentFlags().StringVar(&binPhp, "bin-php", "php", "PHP binary path")

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
