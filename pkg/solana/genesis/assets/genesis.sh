# -*- mode: shell-script -*-
# shellcheck shell=bash

upgradeableLoader=BPFLoaderUpgradeab1e11111111111111111111111
genesis_args=()

fetch-program() {
    local name=$1
    local version=$2
    local address=$3
    local loader=$4

    local so=spl_$name-$version.so

    if [[ $loader == "$upgradeableLoader" ]]; then
        genesis_args+=(--upgradeable-program "$address" "$loader" "$so" none)
    else
        genesis_args+=(--bpf-program "$address" "$loader" "$so")
    fi

    if [[ -r $so ]]; then
        return
    fi

    if [[ -r ~/.cache/solana-spl/$so ]]; then
        cp ~/.cache/solana-spl/"$so" "$so"
    else
        echo "Downloading $name $version"
        local so_name="spl_${name//-/_}.so"
        (
            set -x
            curl -s -S -L --retry 5 --retry-delay 2 --retry-connrefused \
                -o "$so" \
                "https://github.com/solana-labs/solana-program-library/releases/download/$name-v$version/$so_name"
        )

        mkdir -p ~/.cache/solana-spl
        cp "$so" ~/.cache/solana-spl/"$so"
    fi
}

step::000::wait-for-a-stable-environment() {
    cloud-init::wait-for-stable-environment
}

step::005::create-sol-user() {
    create-sol-user
}

step::007::check-for-existing-ledger() {
    if [[ -d $LEDGER_PATH/rocksdb ]]; then
        log::fatal "Ledger directory '$LEDGER_PATH' already appears populated!"
    fi
}

step::010::install-dependencies() {
    svmkit::apt::get install "${PACKAGE_LIST[@]}"
}

step::020::fetch-all-programs() {
    fetch-program token 3.5.0 TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA BPFLoader2111111111111111111111111111111111
    fetch-program token-2022 0.9.0 TokenzQdBNbLqP5VEhdkAS6EPFLC1PHnBqCXEpPxuEb BPFLoaderUpgradeab1e11111111111111111111111
    fetch-program memo 1.0.0 Memo1UhkJRfHyvLMcVucJwxXeuD728EqVDDwQDxFMNo BPFLoader1111111111111111111111111111111111
    fetch-program memo 3.0.0 MemoSq4gqABAXKb96qnH8TysNcWxMyWCqXgDLGmfcHr BPFLoader2111111111111111111111111111111111
    fetch-program associated-token-account 1.1.2 ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL BPFLoader2111111111111111111111111111111111
    fetch-program feature-proposal 1.0.0 Feat1YXHhH6t1juaWF74WLcfv4XoNocjXA6sPWHNgAse BPFLoader2111111111111111111111111111111111
}

step::030::write-primordial-accounts-file() {
    svmkit::sudo cp -f primordial.yaml /home/sol/primordial.yaml
    svmkit::sudo chown sol:sol /home/sol/primordial.yaml
}

step::035::write-validator-accounts-file() {
    if [[ -f validator_accounts.yaml ]]; then
        svmkit::sudo cp -f validator_accounts.yaml /home/sol/validator_accounts.yaml
        svmkit::sudo chown sol:sol /home/sol/validator_accounts.yaml
    fi
}

step::040::execute-solana-genesis() {
    svmkit::sudo -u sol "${GENESIS_ENV[@]}" solana-genesis "${GENESIS_FLAGS[@]}" "${genesis_args[@]}"
}

step::050::create-initial-snapshot() {
    svmkit::sudo -u sol -i agave-ledger-tool create-snapshot ROOT
}
