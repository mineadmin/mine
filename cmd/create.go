package cmd

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mineadmin/mine/internal/downloader"
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
				versions, err := downloader.NewDownloader(language, "", platform).ListVersions()
				if err != nil {
					log.Fatalf("Failed to get versions: %v", err)
				}

				fmt.Println("Available MineAdmin versions:")
				for i, v := range versions {
					fmt.Printf("%d. %s\n", i+1, v)
				}

				var selected int
				fmt.Print("Select version (number): ")
				_, err = fmt.Scanf("%d", &selected)
				if err != nil || selected < 1 || selected > len(versions) {
					log.Fatal("Invalid selection")
				}

				version = versions[selected-1]
			}

			dl := downloader.NewDownloader(language, version, platform)
			fmt.Printf("Creating project %s...\n", projectName)
			fmt.Printf("Language: %s, Version: %s, Platform: %s\n", language, version, platform)

			err := dl.Download(projectName)
			if err != nil {
				log.Fatalf("Failed to create project: %v", err)
			}

			// For PHP projects, collect configuration
			if language == "php" {
				collectConfiguration(projectName)
			}

			fmt.Printf("Successfully created project %s\n", projectName)
		},
	}

	cmd.Flags().StringVarP(&language, "language", "l", "php", "Programming language (php/go/js)")
	cmd.Flags().StringVarP(&version, "version", "v", "latest", "Version of MineAdmin")
	cmd.Flags().StringVarP(&platform, "platform", "p", "swow", "Platform (swow/swoole)")

	return cmd
}

func collectConfiguration(projectDir string) {
	reader := bufio.NewReader(os.Stdin)

	// Database configuration
	fmt.Println("\nDatabase Configuration:")
	dbType := promptSelect(reader, "Database type:", []string{"mysql", "pgsql"})
	dbHost := promptInput(reader, "Database host:", "127.0.0.1")
	dbPort := promptInput(reader, "Database port:", "3306")
	dbName := promptInput(reader, "Database name:", "mineadmin")
	dbUser := promptInput(reader, "Database username:", "root")
	dbPass := promptInput(reader, "Database password:", "root")

	// Redis configuration
	fmt.Println("\nRedis Configuration:")
	redisHost := promptInput(reader, "Redis host:", "127.0.0.1")
	redisPort := promptInput(reader, "Redis port:", "6379")
	redisPass := promptInput(reader, "Redis password (leave empty if none):", "")
	redisDB := promptInput(reader, "Redis database number:", "0")

	// Generate JWT secret
	jwtSecret := generateJWTSecret()

	// Create .env file
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
		log.Fatalf("Failed to create .env file: %v", err)
	}
	fmt.Println("\n.env file created successfully!")
}

func promptInput(reader *bufio.Reader, prompt, defaultValue string) string {
	fmt.Printf("%s [%s]: ", prompt, defaultValue)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue
	}
	return input
}

func promptSelect(reader *bufio.Reader, prompt string, options []string) string {
	fmt.Println(prompt)
	for i, opt := range options {
		fmt.Printf("%d. %s\n", i+1, opt)
	}

	for {
		fmt.Print("Select (number): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		var selected int
		if _, err := fmt.Sscanf(input, "%d", &selected); err != nil || selected < 1 || selected > len(options) {
			fmt.Println("Invalid selection, try again")
			continue
		}

		return options[selected-1]
	}
}

func generateJWTSecret() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Failed to generate JWT secret: %v", err)
	}
	return base64.StdEncoding.EncodeToString(bytes)
}
