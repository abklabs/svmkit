#!/bin/bash

# Constants
RPC_URL="https://zumanet.abklabs.com"
PAYER_WALLET="payer.json"
RECEIVER_WALLET="receiver.json"
AIRDROP_AMOUNT=2
MINT_AMOUNT=100
TRANSFER_AMOUNT=50

# Function to print transaction link
print_tx_link() {
    local signature=$1
    echo -e "\033[1;32m[TX]\033[0m https://zuma.abklabs.com/tx/$signature"
}

# Function to print address link
print_address_link() {
    local address=$1
    echo -e "\033[1;33m[ADDRESS]\033[0m https://zuma.abklabs.com/address/$address"
}

# Function to log info with colors
log_info() {
    local message=$1
    echo -e "\033[1;34m[INFO]\033[0m $message"
}

# Set the RPC endpoint
solana config set --url $RPC_URL

# Create two wallets if they don't already exist
solana-keygen new --outfile $PAYER_WALLET --no-passphrase --force >/dev/null 2>&1
solana-keygen new --outfile $RECEIVER_WALLET --no-passphrase --force >/dev/null 2>&1

payer_pubkey=$(solana-keygen pubkey $PAYER_WALLET)
receiver_pubkey=$(solana-keygen pubkey $RECEIVER_WALLET)

log_info "Payer wallet created with public key:"
print_address_link $payer_pubkey
log_info "Receiver wallet created with public key:"
print_address_link $receiver_pubkey

# Airdrop ZUMA to payer
solana airdrop $AIRDROP_AMOUNT $payer_pubkey

# Create a new token mint and set payer as the mint authority
mint=$(spl-token create-token --mint-authority $payer_pubkey --fee-payer $PAYER_WALLET --output json-compact | jq -r '.commandOutput.address')

# Create an associated token account for payer and receiver
payer_ata_signature=$(spl-token create-account $mint --owner $payer_pubkey --fee-payer $PAYER_WALLET --output json-compact | jq -r '.signature')
receiver_ata_signature=$(spl-token create-account $mint --owner $receiver_pubkey --fee-payer $PAYER_WALLET --output json-compact | jq -r '.signature')

log_info "Created payer associated token account with transaction:"
print_tx_link $payer_ata_signature
log_info "Created receiver associated token account with transaction:"
print_tx_link $receiver_ata_signature

payer_ata_address=$(spl-token address --token $mint --owner $payer_pubkey --verbose --output json-compact | jq -r '.associatedTokenAddress')
receiver_ata_address=$(spl-token address --token $mint --owner $receiver_pubkey --verbose --output json-compact | jq -r '.associatedTokenAddress')

log_info "Payer associated token account address:"
print_address_link $payer_ata_address
log_info "Receiver associated token account address:"
print_address_link $receiver_ata_address

# Mint 100 tokens to payer
mint_to_payer_signature=$(spl-token mint --fee-payer $PAYER_WALLET --mint-authority $PAYER_WALLET --output json-compact $mint $MINT_AMOUNT -- $payer_ata_address | jq -r '.signature')

log_info "Minted 100 tokens to payer with transaction:"
print_tx_link $mint_to_payer_signature

# Check balance of payer
log_info "Balance of payer:"
spl-token balance --owner $payer_pubkey $mint

# Transfer 50 tokens from payer to receiver
payer_to_receiver_signature=$(spl-token transfer --fee-payer $PAYER_WALLET --owner $PAYER_WALLET --output json-compact $mint $TRANSFER_AMOUNT $receiver_ata_address | jq -r '.signature')

log_info "Transfered 50 tokens from payer to receiver with transaction:"
print_tx_link $payer_to_receiver_signature

# Check balances of payer and receiver
log_info "Balance of payer:"
spl-token balance --owner $payer_pubkey $mint

log_info "Balance of receiver:"
spl-token balance --owner $receiver_pubkey $mint
