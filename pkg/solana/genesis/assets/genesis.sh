# -*- mode: shell-script -*-
# shellcheck shell=bash

# shellcheck disable=SC1091
. ./deletion-lib.sh

upgradeableLoader=BPFLoaderUpgradeab1e11111111111111111111111
genesis_args=()

fetch-program() {
    local prefix=$1
    shift
    local name=$1
    shift
    local version=$1
    shift
    local address=$1
    shift
    local loader=$1
    shift
    local url=$1
    shift

    local so=$prefix-$name-$version.so
    local cachedir="$HOME/.cache/solana-$prefix"

    if [[ $loader == "$upgradeableLoader" ]]; then
        genesis_args+=(--upgradeable-program "$address" "$loader" "$so" none)
    else
        genesis_args+=(--bpf-program "$address" "$loader" "$so")
    fi

    if [[ -r $so ]]; then
        return
    fi

    if [[ -r "$cachedir/$so" ]]; then
        cp "$cachedir/$so" "$so"
    else
        echo "Downloading $name $version"
        (
            set -x
            curl -s -S -L --retry 5 --retry-delay 2 --retry-connrefused \
                -o "$so" "$url"

        )

        mkdir -p "$cachedir"
        cp "$so" "$cachedir/$so"
    fi
}

fetch-core-program() {
    local prefix="core-bpf"
    local name="$1"
    local version="$2"
    local so_name="solana_${name//-/_}_program.so"
    local url="https://github.com/solana-program/$name/releases/download/program%40$version/$so_name"

    fetch-program "$prefix" "${@}" "$url"
}

fetch-spl-program() {
    local prefix="spl"
    local name="$1"
    local version="$2"
    local so_name="${prefix}_${name//-/_}.so"
    local url="https://github.com/solana-program/$name/releases/download/program@v$version/$so_name"

    fetch-program "$prefix" "${@}" "$url"
}

step::000::wait-for-a-stable-environment() {
    cloud-init::wait-for-stable-environment
}

step::005::create-sol-user() {
    create-sol-user
}

step::007::check-for-existing-files() {
    if [[ -d $LEDGER_PATH/rocksdb ]]; then
        log::fatal "Ledger directory '$LEDGER_PATH' already appears populated!"
    fi
    deletion::check-create
}

step::010::install-dependencies() {
    svmkit::apt::get install "${PACKAGE_LIST[@]}"
}

step::020::fetch-spl-programs() {
    fetch-spl-program token 3.5.0 TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA BPFLoader2111111111111111111111111111111111
    fetch-spl-program token-2022 8.0.0 TokenzQdBNbLqP5VEhdkAS6EPFLC1PHnBqCXEpPxuEb BPFLoaderUpgradeab1e11111111111111111111111
    fetch-spl-program memo 1.0.0 Memo1UhkJRfHyvLMcVucJwxXeuD728EqVDDwQDxFMNo BPFLoader1111111111111111111111111111111111
    fetch-spl-program memo 3.0.0 MemoSq4gqABAXKb96qnH8TysNcWxMyWCqXgDLGmfcHr BPFLoader2111111111111111111111111111111111
    fetch-spl-program associated-token-account 1.1.2 ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL BPFLoader2111111111111111111111111111111111
    fetch-spl-program feature-proposal 1.0.0 Feat1YXHhH6t1juaWF74WLcfv4XoNocjXA6sPWHNgAse BPFLoader2111111111111111111111111111111111
}

step::025::fetch-core-programs() {
    fetch-core-program address-lookup-table 3.0.0 AddressLookupTab1e1111111111111111111111111 BPFLoaderUpgradeab1e11111111111111111111111
    fetch-core-program config 3.0.0 Config1111111111111111111111111111111111111 BPFLoaderUpgradeab1e11111111111111111111111
    fetch-core-program feature-gate 0.0.1 Feature111111111111111111111111111111111111 BPFLoaderUpgradeab1e11111111111111111111111
    fetch-core-program stake 1.0.0 Stake11111111111111111111111111111111111111 BPFLoaderUpgradeab1e11111111111111111111111
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
