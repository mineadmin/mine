package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mineadmin/mine/internal/downloader"
	"github.com/mineadmin/mine/internal/prompt"
	"github.com/mineadmin/mine/internal/utils"
	"github.com/spf13/cobra"
)

// NewCreateCmd creates and returns the create command
func NewCreateCmd() *cobra.Command {
	var (
		projectName string
		language    string
		version     string
		platform    string
	)

	cmd := &cobra.Command{
		Use:   "create [projectName]",
		Short: "Create a new MineAdmin project",
		Long: `Create a new MineAdmin project with specified language and version.
Example:
  mine create demoProject --language=php --version=v1.0.1 --platform=swow`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			projectName = args[0]

			// For PHP projects, handle version selection if not specified
			if language == "php" && version == "latest" {
				prompt.Info("Fetching available MineAdmin versions...")
				versions, err := downloader.NewDownloader(language, "", platform).ListVersions()
				if err != nil {
					prompt.Error(fmt.Sprintf("Failed to get versions: %v", err))
					os.Exit(1)
				}

				prompt.Info("Select MineAdmin Version")
				_, selectedVersion, err := prompt.Select("Available versions", versions)
				if err != nil {
					prompt.Error(fmt.Sprintf("Version selection failed: %v", err))
					os.Exit(1)
				}
				version = selectedVersion
			}

			dl := downloader.NewDownloader(language, version, platform)
			prompt.Info(fmt.Sprintf("Creating project %s", projectName))
			prompt.Info(fmt.Sprintf("Language: %s, Version: %s, Platform: %s", language, version, platform))

			// Start a spinner for the download process
			spinner := prompt.StartSpinner("Downloading and extracting project files...")
			err := dl.Download(projectName)
			spinner.Stop()

			if err != nil {
				prompt.Error(fmt.Sprintf("Failed to create project: %v", err))
				os.Exit(1)
			}

			// For PHP projects, handle swow platform specific operations
			if language == "php" && platform == "swow" {
				// Check if version > 3.0
				if utils.CompareVersions(version, "3.0") > 0 {
					projectRoot := projectName
					if strings.Contains(projectName, "mineadmin-") {
						projectRoot = filepath.Dir(projectName)
					}

					// Replace files for swow platform
					spinner = prompt.StartSpinner("Configuring project for Swow platform...")

					// Get files from GitHub and replace
					files := []struct {
						srcPath string
						dstPath string
					}{
						{
							srcPath: ".github/ci/hyperf.php",
							dstPath: filepath.Join(projectRoot, "bin", "hyperf.php"),
						},
						{
							srcPath: ".github/ci/server.php",
							dstPath: filepath.Join(projectRoot, "config", "autoload", "server.php"),
						},
						{
							srcPath: ".github/ci/bootstrap.php",
							dstPath: filepath.Join(projectRoot, "tests", "bootstrap.php"),
						},
					}

					for _, file := range files {
						content, err := utils.GetGitHubFileContent("mineadmin/MineAdmin", version, file.srcPath)
						if err != nil {
							spinner.Stop()
							prompt.Error(fmt.Sprintf("Failed to fetch %s from GitHub: %v", file.srcPath, err))
							os.Exit(1)
						}

						// Ensure target directory exists
						if err := os.MkdirAll(filepath.Dir(file.dstPath), 0755); err != nil {
							spinner.Stop()
							prompt.Error(fmt.Sprintf("Failed to create directory for %s: %v", file.dstPath, err))
							os.Exit(1)
						}

						if err := ioutil.WriteFile(file.dstPath, content, 0644); err != nil {
							spinner.Stop()
							prompt.Error(fmt.Sprintf("Failed to write %s: %v", file.dstPath, err))
							os.Exit(1)
						}
					}

					// Modify composer.json
					composerPath := filepath.Join(projectRoot, "composer.json")
					if err := utils.ModifyComposerJSON(composerPath); err != nil {
						spinner.Stop()
						prompt.Error(fmt.Sprintf("Failed to modify composer.json: %v", err))
						os.Exit(1)
					}

					spinner.Stop()
					prompt.Success("Project configured for Swow platform")
				}
			}

			// For PHP projects, collect configuration
			if language == "php" {
				collectConfiguration(projectName)
			}

			prompt.Success(fmt.Sprintf("Successfully created project %s", projectName))
		},
	}

	cmd.Flags().StringVarP(&language, "language", "l", "php", "Programming language (php/go/js)")
	cmd.Flags().StringVarP(&version, "version", "v", "latest", "Version of MineAdmin")
	cmd.Flags().StringVarP(&platform, "platform", "p", "swow", "Platform (swow/swoole)")

	return cmd
}

func collectConfiguration(projectDir string) {
	prompt.Info("Database Configuration")
	spinner := prompt.StartSpinner("Preparing database configuration...")
	_, dbType, err := prompt.Select("Database type", []string{"mysql", "pgsql"})
	spinner.Stop()
	if err != nil {
		prompt.Error(fmt.Sprintf("Database type selection failed: %v", err))
		os.Exit(1)
	}

	dbHost, err := prompt.Input("Database host", "127.0.0.1")
	if err != nil {
		prompt.Error(fmt.Sprintf("Input failed: %v", err))
		os.Exit(1)
	}

	dbPort, err := prompt.Input("Database port", "3306")
	if err != nil {
		prompt.Error(fmt.Sprintf("Input failed: %v", err))
		os.Exit(1)
	}

	dbName, err := prompt.Input("Database name", "mineadmin")
	if err != nil {
		prompt.Error(fmt.Sprintf("Input failed: %v", err))
		os.Exit(1)
	}

	dbUser, err := prompt.Input("Database username", "root")
	if err != nil {
		prompt.Error(fmt.Sprintf("Input failed: %v", err))
		os.Exit(1)
	}

	dbPass, err := prompt.Input("Database password", "root")
	if err != nil {
		prompt.Error(fmt.Sprintf("Input failed: %v", err))
		os.Exit(1)
	}
	prompt.Success("Database configuration completed")

	// Redis configuration
	prompt.Info("Redis Configuration")
	spinner = prompt.StartSpinner("Preparing Redis configuration...")
	redisHost, err := prompt.Input("Redis host", "127.0.0.1")
	spinner.Stop()
	if err != nil {
		prompt.Error(fmt.Sprintf("Input failed: %v", err))
		os.Exit(1)
	}

	redisPort, err := prompt.Input("Redis port", "6379")
	if err != nil {
		prompt.Error(fmt.Sprintf("Input failed: %v", err))
		os.Exit(1)
	}

	redisPass, err := prompt.Input("Redis password (leave empty if none)", "")
	if err != nil {
		prompt.Error(fmt.Sprintf("Input failed: %v", err))
		os.Exit(1)
	}

	redisDB, err := prompt.Input("Redis database number", "0")
	if err != nil {
		prompt.Error(fmt.Sprintf("Input failed: %v", err))
		os.Exit(1)
	}
	prompt.Success("Redis configuration completed")

	// Generate JWT secret
	prompt.Info("Generating security configuration")
	spinner = prompt.StartSpinner("Generating JWT secret...")
	jwtSecret, err := utils.GenerateJwtSecret()
	spinner.Stop()
	if err != nil {
		prompt.Error(fmt.Sprintf("Failed to generate JWT secret: %v", err))
		os.Exit(1)
	}
	prompt.Success("Security configuration completed")

	// Create .env file
	prompt.Info("Creating environment configuration file")
	spinner = prompt.StartSpinner("Writing configuration to .env file...")
	envContent := fmt.Sprintf(`APP_NAME=MineAdmin
APP_ENV=dev
APP_DEBUG=false

DB_DRIVER=%s
DB_HOST=%s
DB_PORT=%s
DB_DATABASE=%s
DB_USERNAME=%s
DB_PASSWORD=%s
DB_CHARSET=utf8mb4
DB_COLLATION=utf8mb4_unicode_ci
DB_PREFIX=

REDIS_HOST=%s
REDIS_AUTH=%s
REDIS_PORT=%s
REDIS_DB=%s

APP_URL=http://127.0.0.1:9501

JWT_SECRET=%s

MINE_ACCESS_TOKEN=(null) # Your MINE_ACCESS_TOKEN
`,
		dbType, dbHost, dbPort, dbName, dbUser, dbPass,
		redisHost, redisPass, redisPort, redisDB,
		jwtSecret)

	// Ensure .env is created in the project root (not in mineadmin-version subdirectory)
	envPath := filepath.Join(projectDir, ".env")

	// If we're in a subdirectory (like mineadmin-version), find the project root
	if strings.Contains(projectDir, "mineadmin-") {
		rootDir := filepath.Dir(projectDir)
		envPath = filepath.Join(rootDir, ".env")
	}
	if err := os.WriteFile(envPath, []byte(envContent), 0644); err != nil {
		spinner.Stop()
		prompt.Error(fmt.Sprintf("Failed to create .env file: %v", err))
		os.Exit(1)
	}
	spinner.Stop()
	prompt.Success("Configuration file created successfully")
}
