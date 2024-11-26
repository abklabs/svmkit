# APT Repository Setup

## Overview

ABK Labs maintains an official APT repository at `apt.abklabs.com` for distributing SVMKit components. More details can be found in the [build](../build/) directory. This repository provides:

- Multiple validator implementations (Agave, Jito, Frankendancer etc.)
- CLI tools and utilities
- Regular updates and security patches
- Seamless integration with Ubuntu/Debian systems

## Available Packages

| Package Name                   | Description                          |
| ------------------------------ | ------------------------------------ |
| `svmkit-agave-validator`       | Agave validator implementation       |
| `svmkit-solana-cli`            | Solana CLI tools                     |
| `svmkit-jito-validator`        | Jito validator implementation        |
| `svmkit-pyth-validator`        | Pyth validator implementation        |
| `svmkit-powerledger-validator` | PowerLedger validator implementation |

## Installation Guide

### 1. System Preparation

First, ensure your system is up-to-date:

```bash
# Update package lists
sudo apt-get update

# Upgrade existing packages (recommended)
sudo apt-get upgrade -y
```

### 2. Install Prerequisites

Install required packages:

```bash
sudo apt-get install -y curl gnupg
```

### 3. Configure Repository

Add the SVMKit repository:

```bash
# Add repository to sources list
echo "deb https://apt.abklabs.com/svmkit dev main" | \
    sudo tee /etc/apt/sources.list.d/svmkit.list

# Import GPG key
curl -s https://apt.abklabs.com/keys/abklabs-archive-dev.asc | \
    sudo apt-key add -

# Update package lists
sudo apt-get update
```

### 4. Install SVMKit Components

**Basic Installation:**
For most users, install these core components:

```bash
sudo apt-get install svmkit-agave-validator svmkit-solana-cli
```

**Available Components**
Choose components based on your needs:

```bash
# Agave validator
sudo apt-get install svmkit-agave-validator

# Solana CLI tools
sudo apt-get install svmkit-solana-cli

# Jito validator
sudo apt-get install svmkit-jito-validator

# Pyth validator
sudo apt-get install svmkit-pyth-validator

# PowerLedger validator
sudo apt-get install svmkit-powerledger-validator
```

## Maintenance

### Updating Packages

```bash
# Update package lists
sudo apt-get update

# Upgrade SVMKit packages
sudo apt-get upgrade 'svmkit-*'
```

### Package Management

```bash
# List installed SVMKit packages
dpkg -l | grep svmkit

# Remove a package
sudo apt-get remove svmkit-agave-validator

# Remove a package and its configuration
sudo apt-get purge svmkit-agave-validator
```

## Troubleshooting

### Common Issues

1. GPG Key Issues

   ```bash
   # Reimport GPG key
   curl -s https://apt.abklabs.com/keys/abklabs-archive-dev.asc | \
       sudo apt-key add -
   ```

2. Repository Access Issues

   ```bash
    # Check repository access
    curl -sI https://apt.abklabs.com/svmkit/dists/dev/Release

    # Verify sources list
    cat /etc/apt/sources.list.d/svmkit.list
   ```

3. Package Conflicts

   ```bash
   # Check package status
   apt-cache policy svmkit-agave-validator

   # Force reinstall if needed
   sudo apt-get install --reinstall svmkit-agave-validator
   ```

## Support

For additional help:

- Join our Telegram
- Create an issue on GitHub
