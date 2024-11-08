#!/bin/bash

# Function to check if a package is installed
check_package_installed() {
    if ! dpkg -l "$1" >/dev/null 2>&1; then
        return 1
    fi
    return 0
}

# Function to check if a repository is configured
check_repository_configured() {
    if ! grep -q "apt.abklabs.com" /etc/apt/sources.list.d/svmkit.list 2>/dev/null; then
        return 1
    fi
    return 0
}

# Function to setup SVMKit requirements
setup_svmkit_requirements() {
    echo "Setting up missing SVMKit requirements..."
    
    # Step 1: Update package lists
    echo "Step 1: Updating package lists..."
    if ! sudo apt-get update; then
        echo "Warning: Initial package list update had issues - continuing anyway..."
    fi
    
    # Step 2: Install prerequisites
    echo "Step 2: Installing prerequisites..."
    if ! sudo apt-get install -y curl gnupg jq ca-certificates; then
        echo "Error: Failed to install prerequisites"
        exit 1
    fi
    
    # Step 3: Add SVMKit repository
    echo "Step 3: Adding SVMKit repository..."
    if ! echo "deb [signed-by=/usr/share/keyrings/abklabs-archive-keyring.gpg trusted=yes] https://apt.abklabs.com/svmkit dev main" | sudo tee /etc/apt/sources.list.d/svmkit.list; then
        echo "Error: Failed to add SVMKit repository"
        exit 1
    fi
    
    # Step 4: Import GPG key using modern method with SSL verification handling
    echo "Step 4: Importing GPG key..."
    if ! curl -fsSL --insecure https://apt.abklabs.com/keys/abklabs-archive-dev.asc > /tmp/abklabs-archive-dev.asc; then
        echo "Error: Failed to download GPG key"
        echo "Please check if https://apt.abklabs.com is accessible"
        exit 1
    fi
    
    if ! sudo gpg --no-default-keyring --keyring /tmp/abklabs.gpg --import /tmp/abklabs-archive-dev.asc; then
        echo "Error: Failed to import GPG key into temporary keyring"
        rm -f /tmp/abklabs-archive-dev.asc /tmp/abklabs.gpg
        exit 1
    fi
    
    if ! sudo gpg --no-default-keyring --keyring /tmp/abklabs.gpg --export | sudo gpg --dearmor -o /usr/share/keyrings/abklabs-archive-keyring.gpg; then
        echo "Error: Failed to export GPG key to system keyring"
        rm -f /tmp/abklabs-archive-dev.asc /tmp/abklabs.gpg
        exit 1
    fi
    
    # Cleanup temporary files
    rm -f /tmp/abklabs-archive-dev.asc /tmp/abklabs.gpg
    
    # Step 5: Configure APT to trust HTTPS for apt.abklabs.com
    echo "Step 5: Configuring APT for HTTPS..."
    if ! echo "Acquire::https::apt.abklabs.com::Verify-Peer \"false\";" | sudo tee /etc/apt/apt.conf.d/99abklabs-ssl; then
        echo "Error: Failed to configure APT HTTPS settings"
        exit 1
    fi
    
    # Step 6: Update package lists again
    echo "Step 6: Updating package lists..."
    if ! sudo apt-get update; then
        echo "Error: Failed to update package lists"
        echo "Please check if the repository is accessible at https://apt.abklabs.com/svmkit"
        exit 1
    fi
    
    # Step 7: Install SVMKit packages
    echo "Step 7: Installing SVMKit packages..."
    if ! sudo apt-get install -y svmkit-agave-validator svmkit-solana-cli; then
        echo "Error: Failed to install SVMKit packages"
        echo "Please verify that the packages are available in the repository"
        echo "You can check by running: apt-cache policy svmkit-agave-validator svmkit-solana-cli"
        exit 1
    fi
    
    echo "SVMKit setup completed successfully!"
}

# Function to verify SSL connection
verify_ssl_connection() {
    local url=$1
    if ! curl -sS --head "$url" >/dev/null 2>&1; then
        echo "Warning: SSL verification failed for $url"
        echo "Would you like to continue with SSL verification disabled? (y/n)"
        read -r response
        if [[ ! "$response" =~ ^[Yy]$ ]]; then
            echo "Setup cancelled."
            exit 1
        fi
        return 1
    fi
    return 0
}

# Function to check SVMKit system requirements
check_svmkit_requirements() {
    local missing_reqs=()
    
    # Check for curl and gnupg
    if ! check_package_installed "curl"; then
        missing_reqs+=("curl")
    fi
    if ! check_package_installed "gnupg"; then
        missing_reqs+=("gnupg")
    fi
    if ! check_package_installed "ca-certificates"; then
        missing_reqs+=("ca-certificates")
    fi
    
    # Check for SVMKit repository configuration
    if ! check_repository_configured; then
        missing_reqs+=("SVMKit repository configuration")
    fi
    
    # Check for required SVMKit packages
    if ! check_package_installed "svmkit-agave-validator"; then
        missing_reqs+=("svmkit-agave-validator")
    fi
    if ! check_package_installed "svmkit-solana-cli"; then
        missing_reqs+=("svmkit-solana-cli")
    fi
    
    if [ ${#missing_reqs[@]} -ne 0 ]; then
        echo "Missing SVMKit requirements: ${missing_reqs[*]}"
        echo "Would you like to automatically install the missing requirements? (y/n)"
        read -r response
        if [[ "$response" =~ ^[Yy]$ ]]; then
            # Verify SSL connection before proceeding
            if ! verify_ssl_connection "https://apt.abklabs.com"; then
                echo "Proceeding with SSL verification disabled..."
            fi
            setup_svmkit_requirements
        else
            echo "Setup cancelled. Please install the requirements manually."
            exit 1
        fi
    fi
}

# Check for required dependencies
check_dependencies() {
    local missing_deps=()
    
    # Check for jq
    if ! command -v jq &> /dev/null; then
        missing_deps+=("jq")
    fi
    
    # Check for solana CLI
    if ! command -v solana &> /dev/null; then
        missing_deps+=("solana CLI tools")
    fi
    
    # Check for spl-token
    if ! command -v spl-token &> /dev/null; then
        missing_deps+=("spl-token")
    fi
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        echo "Missing dependencies: ${missing_deps[*]}"
        echo "Installing missing dependencies..."
        if ! sudo apt-get install -y "${missing_deps[@]}"; then
            echo "Error: Failed to install dependencies"
            exit 1
        fi
    fi
}

# Add error handling for RPC connection
check_rpc_connection() {
    if ! solana cluster-version --url $RPC_URL &> /dev/null; then
        echo "Error: Unable to connect to RPC endpoint: $RPC_URL"
        echo "Please check your internet connection and RPC endpoint."
        exit 1
    fi
}

# Constants
RPC_URL="https://zumanet.abklabs.com"
PAYER_WALLET="payer.json"
RECEIVER_WALLET="receiver.json"
AIRDROP_AMOUNT=2
MINT_AMOUNT=100
TRANSFER_AMOUNT=50

# Check system requirements before proceeding
echo "Checking SVMKit system requirements..."
check_svmkit_requirements

echo "Checking script dependencies..."
check_dependencies

echo "Checking RPC connection..."
check_rpc_connection

# Rest of your existing script...
