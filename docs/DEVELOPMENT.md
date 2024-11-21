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
├── agave/                 # Agave validator implementation
│   ├── assets/           # Validator deployment assets
│   │   └── install.sh    # Validator installation script
│   ├── assets.go         # Asset embedding definitions
│   ├── validator.go      # Core validator logic
│   └── validator_test.go # Validator tests
│
├── genesis/              # Genesis block configuration
│   └── genesis.go        # Genesis creation and management
│
├── runner/               # Deployment and execution system
│   ├── assets/          # Runner deployment assets
│   │   ├── lib.bash     # Common bash utilities
│   │   └── run.sh       # Main runner script
│   ├── assets.go        # Asset embedding definitions
│   ├── envbuilder.go   # Environment variable handling
│   ├── envbuilder_test.go  # Environment builder tests
│   ├── flagbuilder.go      # CLI flag management
│   └── flagbuilder_test.go # Flag builder tests
│   ├── deployer.go      # Remote deployment logic
│   ├── payload.go       # Deployment payload handling
│   └── runner.go        # Main runner implementation
│
├── solana/              # Core Solana functionality
│   ├── assets/         # Solana-specific scripts
│   │   ├── genesis.sh  # Genesis initialization
│   │   ├── stake-account.sh  # Stake account management
│   │   ├── transfer.sh       # Token transfer utilities
│   │   └── vote-account.sh   # Vote account management
│   ├── assets.go       # Asset embedding definitions
│   ├── cli.go          # CLI interface implementation
│   ├── cli_test.go     # CLI tests
│   ├── env.go          # Environment configuration
│   ├── genesis.go      # Genesis block management
│   ├── stakeaccount.go # Stake account operations
│   ├── transfer.go     # Transfer operations
│   └── voteaccount.go  # Vote account operations
│
└── validator/          # Generic validator interface
    └── validator.go    # Validator interface definitions
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
