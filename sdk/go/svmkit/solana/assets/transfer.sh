# -*- mode: shell-script -*-
# shellcheck shell=bash

umask 077

transfer-create() {
    local args=("${SOLANA_CLI_TXN_FLAGS[@]}")

    if [[ -v ALLOW_UNFUNDED_RECIPIENT ]]; then
        args+=(--allow-unfunded-recipient)
    fi

    solana transfer "${args[@]}" "$RECIPIENT_PUBKEY" "$AMOUNT"
}

case "$TRANSFER_ACTION" in
CREATE)
    transfer-create
    ;;
*)
    log::fatal "unknown action provided '$TRANSFER_ACTION'!"
    ;;
esac
