# MineAdmin CLI Tool

A powerful command line tool for downloading and managing MineAdmin projects.

## Features
- One-click project creation
- Support for multiple programming languages (PHP, Go, JavaScript)
- Support for different runtime platforms (Swow, Swoole)
- Interactive configuration for database and Redis connections
- Automatic generation of security keys and environment config files
- Automatic dependency installation and database migrations
- Rich command line interaction experience

## Installation
```bash
go install github.com/mineadmin/mine@latest
```

## Usage
### Create a new project
```bash
mine create <project_name> [--language=<language>] [--version=<version>] [--platform=<platform>]
```

Defaults:
- language: php
- version: latest
- platform: swow

### List available versions
```bash
mine select-versions --language=<language>
```

## Supported Languages
- **PHP** (fully supported)
  - Downloads from GitHub releases
  - Interactive database and Redis configuration
  - Automatic .env file generation
  - Automatic dependency installation and migrations
- **Go** (coming soon)
- **JavaScript** (coming soon)

## Supported Platforms
- **swow** (default)
  - High-performance coroutine runtime
  - Auto-configures project for Swow
- **swoole**
  - Traditional PHP coroutine runtime
- **none** (for non-PHP projects)

## Project Structure
```
.
├── cmd/                # Command implementations
│   ├── create.go       # Create project command
│   ├── root.go         # Root command and main entry
│   └── select_versions.go # Version selection command
├── internal/           # Internal packages
│   ├── downloader/     # Core download functionality
│   │   └── downloader.go
│   ├── prompt/         # CLI interaction functionality
│   │   └── prompt.go
│   └── utils/          # Utility functions
│       └── utils.go
├── main.go             # CLI entry point
├── go.mod              # Go module definition
└── go.sum              # Dependency checksums
```

## Core Components
1. **Downloader** (internal/downloader/downloader.go)
   - Handles downloading MineAdmin projects from GitHub
   - Supports:
     - Downloading PHP projects from GitHub releases
     - Listing versions via GitHub API
     - File extraction and project setup

2. **Prompt** (internal/prompt/prompt.go)
   - Provides rich command line interaction experience
   - Features:
     - User input collection and validation
     - Selection lists
     - Colored output
     - Progress indicators

3. **Utils** (internal/utils/utils.go)
   - Provides various utility functions
   - Features:
     - Random key generation
     - Command and extension checking
     - Semantic version comparison
     - File operations
     - Command execution

4. **CLI Commands**
   - Root command (cmd/root.go): Defines main CLI structure
   - Create command (cmd/create.go): Handles project creation
   - Select-versions command (cmd/select_versions.go): Lists available versions

## Building from Source
```bash
# Clone repository
git clone https://github.com/mineadmin/mine.git
cd mine

# Build binary
go build -o mine

# Install
go install
```

## Dependencies
This project uses these main dependencies:
- github.com/spf13/cobra v1.9.1 (CLI framework)
- github.com/briandowns/spinner v1.23.0 (terminal loading animation)
- github.com/fatih/color v1.16.0 (terminal colors)
- github.com/manifoldco/promptui v0.9.0 (interactive prompts)

## About MineAdmin
MineAdmin is a high-performance PHP backend management system that supports Swoole and Swow coroutine runtimes. It provides rich features including permission management, system monitoring, code generators, etc.

- Official website: https://www.mineadmin.com
- GitHub: https://github.com/mineadmin/mineadmin

## Contributing
Welcome to contribute code, report issues or suggest new features. Please follow these steps:

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add some amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Create Pull Request

## Project Links
- GitHub: https://github.com/mineadmin/mine
- MineAdmin: https://github.com/mineadmin/mineadmin

## License
MIT License

Copyright (c) [year] [fullname]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.