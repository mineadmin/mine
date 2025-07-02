# MineAdmin CLI Tool

A command line tool for downloading and managing MineAdmin projects.

## Installation

```bash
go install github.com/mineadmin/mine
```

## Usage

### Create a new project

```bash
mine create <projectName> [--language=<language>] [--version=<version>] [--platform=<platform>]
```

Defaults:
- language: php
- version: latest
- platform: swow

Example:
```bash
mine create demoProject  # Uses all defaults
mine create demoProject --version=v1.0.1  # Override version only
```

### List available versions

```bash
mine select-versions --language=<language>
```

Example:
```bash
mine select-versions --language=php
```

## Advanced Usage

### Supported Languages

- PHP (fully supported)
  - Downloads from GitHub releases
  - Supports interactive configuration for database and Redis
  - Auto-generates .env file
- Go (coming soon)
- JavaScript (coming soon)

### Supported Platforms

- swow (default)
- swoole
- none (for non-PHP projects)

### Command Options Details

#### create command

### Dependencies
This project uses the following main dependencies:
- github.com/spf13/cobra v1.9.1 (CLI framework)
- github.com/briandowns/spinner v1.23.0 (terminal spinner)
- github.com/fatih/color v1.16.0 (terminal colors)
- github.com/manifoldco/promptui v0.9.0 (interactive prompts)

- `--language`: Specify project language (php/go/js)
- `--version`: Specify MineAdmin version (default: latest)
- `--platform`: Runtime platform (swow/swoole/none)

#### select-versions command
- `--language`: Query versions for specific language
  - For PHP: Fetches from GitHub API
  - Others: Returns mock data (to be implemented)

## Development

### Project Structure

```
.
├── cmd/                # Command implementations
│   ├── create.go       # Create project command
│   ├── root.go         # Root command and main entry
│   └── select_versions.go # Version selection command
├── internal/
│   └── downloader/     # Core download functionality
│       └── downloader.go
├── main.go             # CLI entry point
├── go.mod              # Go module definition
└── go.sum              # Dependency checksums
```

### Core Components

1. **Downloader** (internal/downloader/downloader.go)
   - Handles downloading MineAdmin projects from GitHub
   - Supports:
     - PHP project downloads from GitHub releases
     - Version listing via GitHub API
     - File extraction and project setup

2. **CLI Commands**
   - Root command (cmd/root.go): Defines the main CLI structure
   - Create command (cmd/create.go): Handles project creation
   - Select-versions command (cmd/select_versions.go): Lists available versions

### Building from Source

```bash
# Clone the repository
git clone https://github.com/mineadmin/mine.git
cd mine

# Build the binary
go build -o mine

# Install
go install
```

## Project Links

- GitHub: https://github.com/mineadmin/mine
- MineAdmin: https://github.com/mineadmin/mineadmin