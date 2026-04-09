---
title: Essential CLI Tool Guide for the AI Era
icon: /icons/streamline-ultimate-blog-blogger-logo.svg
---

In the AI era, Command Line Interface (CLI) tools are experiencing a **renaissance**. Why? Because **CLI tools are naturally AI-Agent friendly**!

Compared to graphical interfaces, CLI tools have structured input/output, clear help documentation, and predictable behavior patterns—these characteristics make it much easier for AI Agents to understand, learn, and automate CLI tool usage.

But have you ever struggled with creating a fully-functional CLI tool:

- Tedious command-line argument parsing?
- Chaotic configuration management?
- Difficult-to-maintain help documentation?
- Want to extend with plugins but don't know where to start?

---

💡 **First, let's introduce the jzero template market**

jzero provides a rich collection of **official templates** and **third-party templates** to help you quickly build various types of projects:

**🚀 Built-in Templates**:
- **RPC Template**: High-performance gRPC microservice framework based on Protocol Buffers
- **API Template**: Lightweight RESTful API service framework based on API description language
- **Gateway Template**: High-performance API gateway supporting both gRPC and HTTP protocols

**📦 Official External Templates**:
- **CLI Template**: Command-line application template with common CLI patterns (today's star!)
- **API Template**: API template optimized for Vercel deployment
- **Documentation Template**: Documentation site template using VuePress Hope theme

**🌍 Third-party Templates**:
- Contributions welcome! Share your own template to help more developers start projects quickly!

Visit **[jzero Template Market](https://templates.jzero.io/)** for more template information and usage guides.

![](https://oss.jaronnie.com/image-20260409190255415.png)

---

Today, we'll introduce how to use the **jzero CLI template** to rapidly build professional command-line tools!

![](https://oss.jaronnie.com/image-20260409190335125.png)

---

## Why Choose jzero CLI Template?

The jzero CLI template is built on the industry-standard **Cobra framework**, providing an out-of-the-box project structure and best-practice configurations. Compared to building from scratch, using the jzero CLI template enables you to:

✅ **Quick Start**: Generate a complete project structure with one click, no tedious configuration needed  
✅ **Unified Standards**: Follow industry standards with clear, understandable command structure  
✅ **Feature-Complete**: Built-in enterprise-grade features like configuration management, plugin system, debug mode  
✅ **Easy Extension**: Plugin architecture for easily adding new features  
✅ **AI-Friendly**: Perfect cooperation with Claude, GPT, and other AI tools to boost development efficiency  

---

## Quick Start: Create Your First CLI Tool in 3 Minutes

```bash
# 1. Install jzero (if not already installed)
go install github.com/jzero-io/jzero/cmd/jzero@latest

# 2. Create a new CLI project
jzero new mycli --branch cli

cd mycli

# 3. Install dependencies and build
go mod tidy
go build

# 4. Test run
./mycli version
```

Output example:
```
mycli version 1.0.0 darwin/amd64
Go version go1.21.0
Git commit abc123
Build date 2024-01-01 12:00:00 +0000 UTC
```

It's that simple! You now have a fully-functional CLI tool framework.

The project structure is as follows:

```
mycli/
├── main.go                    # Entry file
├── internal/
│   ├── command/              # Command implementations
│   │   └── version/          # Version command
│   │       └── version.go
│   └── config/               # Configuration management
│       └── config.go
├── go.mod
└── go.sum
```

---

## Core Concepts: Three-Layer Command Structure

The jzero CLI template is based on the Cobra framework, using a clear **three-layer command structure**:

```
Root Command
├── Command
│   └── Sub Command
```

### 1. Root Command

The **Root Command** is the entry point of the CLI tool, defining basic information, global configuration, and top-level commands.

```go
// main.go
var rootCmd = &cobra.Command{
    Use:   "mycli",
    Short: "My CLI tool",
    Long:  `A powerful command-line tool to help improve your productivity`,
}
```

Root command features:
- ✅ Define global flags (like `--debug`, `--config`)
- ✅ Register top-level commands
- ✅ Provide overall help information for the tool

### 2. Command

**Commands** are direct subcommands under the root command, representing major functional modules.

```bash
mycli version      # Version command - display version information
mycli config       # Config command - manage configuration
mycli plugin       # Plugin command - manage plugins
mycli server       # Server command - start server
```

Command characteristics:
- ✅ Directly mounted under root command
- ✅ Can have independent flags and arguments
- ✅ Can contain subcommands, forming a command tree

### 3. Sub Command

**Sub Commands** are the next level under commands, used to implement more granular functionality.

```bash
# Subcommands of config command
mycli config init      # Initialize configuration
mycli config set       # Set configuration item
mycli config get       # Get configuration item
mycli config list      # List all configurations

# Subcommands of plugin command
mycli plugin install   # Install plugin
mycli plugin remove    # Remove plugin
mycli plugin list      # List plugins
mycli plugin update    # Update plugin
```

Subcommand advantages:
- ✅ Modular functionality with clear logic
- ✅ Support multi-level nesting (e.g., `mycli config database connect`)
- ✅ Each subcommand can be developed and maintained independently

### Command Examples Comparison

```bash
# Root Command
mycli                    # Execute root command

# Command (Level 1)
mycli config             # Execute config command
mycli plugin             # Execute plugin command

# Sub Command (Level 2)
mycli config init        # Execute config init subcommand
mycli plugin install     # Execute plugin install subcommand

# Deeper subcommands (Level 3)
mycli server start       # Start server
mycli server stop        # Stop server
mycli server status      # View server status
```

### Types of Flags

Flags configure command behavior, divided into three types:

**Local Flags** - Only valid for current command:
```go
Cmd.Flags().StringP("output", "o", "", "Output file")
```

**Persistent Flags** - Valid for current command and all its subcommands:
```go
rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output mode")
```

**Global Flags** - Valid for all commands:
```bash
mycli --debug          # Enable debug mode
mycli --config file.yaml  # Specify configuration file
```

---

## Configuration Management: Flexible Multi-Layer Configuration

The jzero CLI template provides a powerful configuration management system supporting three configuration methods:

### Configuration Priority
```
Environment variables > Command-line flags > Configuration file
```

### 1. Configuration File

Default configuration file location: `~/.mycli.yaml`

```yaml
# Debug mode
debug: false

# Debug sleep time (seconds)
debug-sleep-time: 0

# Greeting configuration
greet:
  name: World
```

### 2. Environment Variable Configuration


The jzero CLI template automatically maps environment variables to configuration fields:

```bash
# Set environment variables directly, auto-mapped to config
export MYCLI_DEBUG=true
export MYCLI_DEBUG_SLEEP_TIME=5
export MYCLI_GREET_NAME="Alice"
```

**Environment Variable Naming Rules**:

Format: `{APP_PREFIX}_{CONFIG_PATH}`

- `{APP_PREFIX}`: App name prefix (uppercase), e.g., `MYCLI`, `JZERO`
- `{CONFIG_PATH}`: Configuration path with `.` and `-` replaced by `_`

Mapping examples:
| Config Field | Environment Variable |
|-------------|---------------------|
| `config.C.Debug` | `MYCLI_DEBUG` |
| `config.C.DebugSleepTime` | `MYCLI_DEBUG_SLEEP_TIME` |
| `config.C.Greet.Name` | `MYCLI_GREET_NAME` |

### 3. Command-line Flags

```bash
# Override configuration via command line
./mycli --debug
./mycli --config custom.yaml
```

---

## Unified Configuration Management
The jzero CLI template provides a unified configuration management system where all configuration is managed through `internal/config/config.go`. It supports three methods: configuration files, environment variables, and command-line flags, with automatic priority-based loading.

Configuration priority: **Command-line flags > Environment variables > Configuration file**

### Advantages of Unified Configuration

✅ **Automatic Priority Handling**: Viper automatically handles configuration priority
✅ **Type Safe**: Use Go structs with mapstructure tags
✅ **Environment Variable Support**: Automatically maps environment variables to config fields
✅ **Flexible Extension**: Adding new configuration fields is very simple
✅ **Global Access**: Access from anywhere via `config.C` global variable
✅ **Command-Specific Config**: Supports independent configuration for different commands
---

## Adding Custom Commands: Four Steps Complete

Adding a new command requires four steps, following the actual development workflow:

### Step 1: Create Command File

Create a new directory and command file under `internal/command/`:

```go
// internal/command/greet/greet.go
package greet

import (
    "fmt"
    "mycli/internal/config"
    "github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
    Use:   "greet",
    Short: "Greeting command",
    Long:  `Send friendly greetings to users`,
    Run: func(cmd *cobra.Command, args []string) {
        // Get name from unified configuration (mapped to greet.name)
        name := config.C.Greet.Name
        fmt.Printf("Hello, %s!\n", name)
    },
}
```

### Step 2: Define Command Flags

Define flags in the command's `init()` function:

```go
func init() {
    // Define short and long flags with default value
    Cmd.Flags().StringP("name", "n", "World", "Specify the name to greet")
}
```

**Important Notes**:
- Flag names use lowercase (e.g., `name`)
- Automatically maps to `config.C.Greet.Name`
- Short flag `-n` is a shorthand for long flag `--name`

### Step 3: Add Configuration Fields

Add corresponding configuration structures in `internal/config/config.go`:

```go
type Config struct {
    Debug bool `mapstructure:"debug"`
    DebugSleepTime int `mapstructure:"debug-sleep-time"`

    // Add greet command configuration
    Greet GreetConfig `mapstructure:"greet"`
}

// GreetConfig configuration for greet command
type GreetConfig struct {
    Name string `mapstructure:"name"`
}
```

**Configuration Mapping Rules**:
- Command flag `name` → Config field `Greet.Name`
- Config file `greet.name` → `config.C.Greet.Name`
- Environment variable `MYCLI_GREET_NAME` → `config.C.Greet.Name`

### Step 4: Register Command

Import and register the command in `main.go`:

```go
import (
    "mycli/internal/command/greet"
    // other imports...
)

func init() {
    rootCmd.AddCommand(greet.Cmd)
}
```

### Test Usage

```bash
go build

# Method 1: Use default value
./mycli greet
# Output: Hello, World!

# Method 2: Use command-line flag
./mycli greet --name Alice
# Output: Hello, Alice!

# Method 3: Use configuration file
echo "greet:" >> ~/.mycli.yaml
echo "  name: Bob" >> ~/.mycli.yaml
./mycli greet
# Output: Hello, Bob!

# Method 4: Use environment variable
export MYCLI_GREET_NAME="Charlie"
./mycli greet
# Output: Hello, Charlie!
```

---

## Debug Mode: Developer's Best Friend

The jzero CLI template includes comprehensive debug support:

### Three Ways to Enable Debugging

**Method 1: Configuration File**
```yaml
# ~/.mycli.yaml
debug: true
debug-sleep-time: 5  # Debug sleep time (seconds)
```

**Method 2: Environment Variables**
```bash
export MYCLI_DEBUG=true
./mycli
```

**Method 3: Command-line Flags**
```bash
./mycli --debug
./mycli --debug --debug-sleep-time 5
```

### Debug Features

- **Verbose Logging**: Display detailed execution process and intermediate states
- **Sleep Time Control**: Pause between key steps for observation
- **Error Stack Traces**: Clear error information and call stacks

---

## Plugin System: Infinite Extension Possibilities

The jzero CLI template supports a powerful plugin system, allowing your tool to dynamically extend functionality.

### Plugin Naming Rules

Plugin executable files are prefixed with `YOUR_APP-`, for example:
- `mycli-git`
- `mycli-docker`
- `mycli-deploy`

### Plugin Auto-Discovery

The system automatically searches for executable files starting with `mycli-` from PATH.

### Plugin Usage Examples

```bash
# Install plugin to PATH
sudo cp mycli-git /usr/local/bin/

# Use plugin directly
./mycli git status
./mycli docker build
```
## Distribution: Using GoReleaser and GitHub Workflows

After development is complete, how do you conveniently distribute your CLI tool? **GoReleaser** combined with **GitHub Actions** enables automated building and releasing!

### GoReleaser Configuration

Create `.goreleaser.yaml` in the project root:

```yaml
# .goreleaser.yaml
version: 2

before:
  hooks:
    - go mod tidy

builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
      - windows
      - darwin
    goarch:
      - amd64
      - arm64
    dir: ./cmd/mycli
    id: mycli
    binary: mycli
    ldflags:
      - -s -w
      - -X "main.version={{.Version}}"
      - -X "main.commit={{.Commit}}"
      - -X "main.date={{.Date}}"

archives:
  - formats: [tar.gz]
    name_template: >-
      {{ .ProjectName }}_
      {{- title .Os }}_
      {{- if eq .Arch "amd64" }}x86_64
      {{- else if eq .Arch "386" }}i386
      {{- else }}{{ .Arch }}{{ end }}
      {{- if .Arm }}v{{ .Arm }}{{ end }}
    format_overrides:
      - goos: windows
        formats: [zip]

checksum:
  name_template: "checksums.txt"

changelog:
  sort: asc
  filters:
    exclude:
      - "^docs:"
      - "^test:"

force_token: github
```

### GitHub Actions Configuration

Create `.github/workflows/release.yml`:

> **Note**: Before using, you need to configure `ACCESS_TOKEN` in GitHub repository secrets (Settings → Secrets and variables → Actions → New repository secret) with write access to repository. You can generate a Personal Access Token (PAT) for this purpose.
>
> ![](https://oss.jaronnie.com/image-20260410113251665.png)

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
      - name: Set up QEMU
        uses: docker/setup-qemu-action@v4

      - name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Get version
        id: get_version
        run: echo ::set-output name=VERSION::${GITHUB_REF/refs\/tags\//}

      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v5
        with:
          distribution: goreleaser
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.ACCESS_TOKEN }}

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v4

      - name: Log in to Docker Hub
        uses: docker/login-action@v4
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.ACCESS_TOKEN }}

      - name: Docker build and push
        uses: docker/build-push-action@v6
        with:
          registry: ghcr.io
          context: .
          platforms: linux/amd64,linux/arm64
          push: true
          tags: |
            ghcr.io/${{ github.repository }}:latest
            ghcr.io/${{ github.repository }}:${{ steps.get_version.outputs.VERSION }}
```

### Dockerfile Configuration

Create `Dockerfile` for containerized deployment:

```dockerfile
FROM alpine:latest

ENV CGO_ENABLED=0

LABEL \
    org.opencontainers.image.title="mycli" \
    org.opencontainers.image.description="My CLI tool" \
    org.opencontainers.image.url="https://github.com/yourusername/mycli" \
    org.opencontainers.image.documentation="https://github.com/yourusername/mycli#readme" \
    org.opencontainers.image.source="https://github.com/yourusername/mycli" \
    org.opencontainers.image.licenses="MIT" \
    maintainer="your-name <your-email@example.com>"

WORKDIR /app

COPY dist/mycli_linux_amd64_v1/mycli /dist/mycli_linux_amd64/mycli
COPY dist/mycli_linux_arm64_v8.0/mycli /dist/mycli_linux_arm64/mycli

RUN if [ $(go env GOARCH) = "amd64" ]; then \
      cp /dist/mycli_linux_amd64/mycli /usr/local/bin/mycli; \
    elif [ $(go env GOARCH) = "arm64" ]; then \
      cp /dist/mycli_linux_arm64/mycli /usr/local/bin/mycli; \
    fi

RUN apk update --no-cache \
    && apk add --no-cache tzdata ca-certificates \
    && rm -rf /dist

ENTRYPOINT ["mycli"]
```

### Release Process

**1. Create and push a tag**:

```bash
# Create version tag
git tag v1.0.0

# Push tag to GitHub
git push origin v1.0.0
```

**2. GitHub Actions executes automatically**:

- Triggers `.github/workflows/release.yml`
- GoReleaser automatically builds multi-platform binaries
- Generates GitHub Release
- Uploads build artifacts and checksums
- Builds and pushes Docker images to GHCR

**3. User installation methods**:

**Method 1: Direct binary download**

```bash
# Download from GitHub Releases
wget https://github.com/yourusername/mycli/releases/download/v1.0.0/mycli_1.0.0_linux_amd64.tar.gz

tar -xzf mycli_1.0.0_linux_amd64.tar.gz
sudo mv mycli /usr/local/bin/
```

**Method 2: Install using go install command**

```bash
# Install specific version directly
go install github.com/yourusername/mycli@v1.0.0

# Install latest version
go install github.com/yourusername/mycli@latest
```

**Method 3: Using Docker Image**

```bash
# Run directly
docker run --rm ghcr.io/yourusername/mycli:latest version

# Create an alias for convenience
alias mycli='docker run --rm -v $(pwd):/app -w /app ghcr.io/yourusername/mycli:latest'

# Then use it like a local installation
mycli version
mycli --help
```

Or create a shell script `mycli-docker.sh`:

```bash
#!/bin/bash
docker run --rm -v "$(pwd)":/app -w /app ghcr.io/yourusername/mycli:latest "$@"
```

Then add it to your PATH:

```bash
chmod +x mycli-docker.sh
sudo mv mycli-docker.sh /usr/local/bin/mycli
```

### Advantages of Automated Distribution Process

✅ **Fully Automated**: End-to-end automation from code tag to release completion
✅ **Multi-platform Support**: Build once, support Linux, macOS, Windows with multiple architectures
✅ **Multiple Distribution Formats**: Binary files, Docker images, go install - multiple installation methods
✅ **Version Management**: Automatically inject version info, manage releases through tags
✅ **Security Assurance**: Automatically generate SHA256 checksums to ensure download integrity
✅ **Container Deployment**: Automatically build and push multi-arch Docker images to GHCR

---

## Complete Example

To help you better understand how to build command-line tools using the jzero CLI template, we've created a complete demonstration project.

**Project URL**: [https://github.com/jaronnie/mycli](https://github.com/jaronnie/mycli)

### Project Features

This demo project is created entirely following the workflow in this guide and includes:

- ✅ **greet command**: Demonstrates how to add custom commands
- ✅ **Unified configuration management**: Shows usage of config files, environment variables, and command-line flags
- ✅ **GoReleaser configuration**: Complete cross-platform build configuration
- ✅ **GitHub Actions workflow**: Automated release process
- ✅ **Dockerfile**: Multi-architecture containerization support
- ✅ **Complete documentation**: Detailed README and usage instructions

### Quick Try

```bash
# Clone the project
git clone https://github.com/jaronnie/mycli.git
cd mycli

# Install dependencies and build
go mod tidy
go build

# Test run
./mycli version
./mycli greet
./mycli greet --name Alice
```

### Project Structure

```
mycli/
├── main.go                    # Entry file
├── internal/
│   ├── command/              # Command implementations
│   │   ├── version/          # Version command
│   │   └── greet/            # Greet command (custom)
│   └── config/               # Configuration management
│       └── config.go
├── Dockerfile                 # Docker configuration
├── .goreleaser.yaml          # GoReleaser configuration
├── .github/workflows/
│   └── release.yml           # GitHub Actions
├── go.mod
├── go.sum
└── README.md                 # Project documentation
```

### Use as Project Template

You can directly use the mycli project as a base to develop your own CLI tool:

```bash
# Fork or clone the project
git clone https://github.com/jaronnie/mycli.git your-cli
cd your-cli

# Modify configuration
# - Edit app name and description in main.go
# - Change module path in go.mod
# - Add your own commands
# - Update README.md

# Start developing!
```

This project demonstrates best practices for the jzero CLI template and serves as an excellent starting point for learning CLI tool development.

![](https://oss.jaronnie.com/image-20260410115131704.png)

![](https://oss.jaronnie.com/image-20260410115221574.png)

---

## Related Resources

- **jzero GitHub**: [https://github.com/jzero-io/jzero](https://github.com/jzero-io/jzero)
- **jzero Documentation**: [https://docs.jzero.io](https://docs.jzero.io)
- **CLI Template**: [https://templates.jzero.io/external/cli/](https://templates.jzero.io/external/cli/)
- **Cobra Documentation**: [https://github.com/spf13/cobra](https://github.com/spf13/cobra)
- **Viper Documentation**: [https://github.com/spf13/viper](https://github.com/spf13/viper)

---

**Let jzero CLI template be your capable assistant in the AI era!** 🚀

**Find it useful? Please give jzero a ⭐ Star to support our continued improvement!**