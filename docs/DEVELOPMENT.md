# Development Guide

## Overview

This guide provides comprehensive information for developers who want to contribute to or build with SVMKit. Whether you're looking to submit a PR, build your own fork, or integrate SVMKit into your project, you'll find everything you need here.

## Requirements

### Core Tools

- [Go 1.22+](https://golang.org/dl/) - Primary development language
- [golangci-lint](https://golangci-lint.run/install) - Code linting and static analysis
- [make](https://www.gnu.org/software/make/) - Build automation

### Installation

```bash
# macOS (using Homebrew)
brew install go golangci-lint make

# Ubuntu/Debian
sudo apt-get update
sudo apt-get install golang-1.22 make
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Verify installations
go version
~/go/bin/golangci-lint --version
make --version
```

## Project Structure

```sh
svmkit/
├── agave/       # Agave validator implementation
│   ├── assets/  # Validator deployment assets
│
├── runner/      # Deployment and execution system
│   ├── assets/  # Runner deployment assets
│
└── solana/      # Core Solana functionality
    └── assets/  # Solana-specific scripts
```

## Development Workflow

### 1. Setting Up Development Environment

#### Clone repository

```sh
git clone https://github.com/abklabs/svmkit.git
cd svmkit
```

### 2. Testing

```sh
# Run all tests
make test

# Run specific test
go test ./pkg/agave -run TestValidatorEnv

```

## Common Development Tasks

### Adding a New Validator Fork

1. Add fork configuration in `pkg/agave/validator.go`:

   ```go
   Copyconst (
       VariantNewFork Variant = "newfork"
   )
   ```

2. Implement fork-specific logic in `pkg/agave/assets/`
3. Update build scripts in `build/`
4. Add tests
5. Update documentation
