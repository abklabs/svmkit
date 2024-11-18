# -*- mode: shell-script -*-
# shellcheck shell=bash

umask 077

transfer-create() {
    local payer_keypair args

    payer_keypair=$(temp::file)
    echo "$PAYER_KEYPAIR" >"$payer_keypair"

    args=()

    if [[ -v ALLOW_UNFUNDED_RECIPIENT ]]; then
        args+=(--allow-unfunded-recipient)
    fi

    solana -k "$payer_keypair" transfer "${args[@]}" "$RECIPIENT_PUBKEY" "$AMOUNT"
}

case "$TRANSFER_ACTION" in
CREATE)
    transfer-create
    ;;
*)
    log::fatal "unknown action provided '$TRANSFER_ACTION'!"
    ;;
esac
