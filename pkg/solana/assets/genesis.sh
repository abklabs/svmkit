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
    if [[ -d $LEDGER_PATH/rocksdb ]] ; then
	log::fatal "Ledger directory '$LEDGER_PATH' already appears populated!"
    fi
}

step::010::install-dependencies() {
    apt::setup-abk-apt-source

    if [[ -v PACKAGE_VERSION ]]; then
        $APT install bzip2 "svmkit-solana-genesis=$PACKAGE_VERSION" "svmkit-solana-cli=$PACKAGE_VERSION" "svmkit-agave-ledger-tool=$PACKAGE_VERSION"
    else
        $APT install bzip2 svmkit-solana-genesis svmkit-solana-cli svmkit-agave-ledger-tool
    fi
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
    if [[ -z "$PRIMORDIAL_PUBKEYS" || -z "$PRIMORDIAL_LAMPORTS" ]]; then
        log::info "PRIMORDIAL_PUBKEYS or PRIMORDIAL_LAMPORTS variable is not set or empty. Primordial file will be empty."
        $SUDO -u sol tee /home/sol/primordial.yaml </dev/null >/dev/null
        return 0
    fi

    local pubkeys=(${PRIMORDIAL_PUBKEYS//,/ })
    local lamports=(${PRIMORDIAL_LAMPORTS//,/ })

    if [[ ${#pubkeys[@]} -ne ${#lamports[@]} ]]; then
        log::error "The number of pubkeys and lamports entries do not match."
        return 1
    fi

    for i in "${!pubkeys[@]}"; do
        local pubkey=${pubkeys[$i]}
        local lamport=${lamports[$i]}
        cat <<EOF | $SUDO -u sol tee -a /home/sol/primordial.yaml >/dev/null
$pubkey:
  balance: $lamport
  owner: 11111111111111111111111111111111
  executable: false
  data:
EOF
    done
}

step::040::execute-solana-genesis() {
    $SUDO -u sol solana-genesis \
        --ledger $LEDGER_PATH \
        --bootstrap-validator \
        $IDENTITY_PUBKEY \
        $VOTE_PUBKEY \
        $STAKE_PUBKEY \
        --faucet-pubkey $FAUCET_PUBKEY \
        --faucet-lamports $FAUCET_LAMPORTS \
        --target-lamports-per-signature $TARGET_LAMPORTS_PER_SIGNATURE \
        --inflation $INFLATION \
        --lamports-per-byte-year $LAMPORTS_PER_BYTE_YEAR \
        --slots-per-epoch $SLOT_PER_EPOCH \
        --cluster-type $CLUSTER_TYPE \
        --primordial-accounts-file /home/sol/primordial.yaml \
        "${genesis_args[@]}"
}

step::050::create-initial-snapshot() {
    $SUDO -u sol -i agave-ledger-tool create-snapshot --ignore-ulimit-nofile-error ROOT
}
