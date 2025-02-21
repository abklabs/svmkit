# -*- mode: shell-script -*-
# shellcheck shell=bash

umask 077

stake-account-create () {
    solana create-stake-account "${SOLANA_CLI_TXN_FLAGS[@]}" stake_account.json "$STAKE_AMOUNT"
    solana delegate-stake "${SOLANA_CLI_TXN_FLAGS[@]}" stake_account.json vote_account.json
}

stake-account-delete () {
    [[ -v FORCE_DELETE ]] || return 0

    local args=()
    if [[ -v ADD_WITHDRAW_AUTHORITY ]]; then
      args+=(--withdraw-authority withdraw_authority.json)
    fi

    solana withdraw "${SOLANA_CLI_TXN_FLAGS[@]}" "${args[@]}" stake_account.json $WITHDRAW_PUBKEY "$STAKE_AMOUNT"
}

case "$STAKE_ACCOUNT_ACTION" in
    CREATE)
	stake-account-create
	;;
    DELETE)
	stake-account-delete
	;;
    *)
	log::fatal "unknown action provided '$STAKE_ACCOUNT_ACTION'!"
    ;;
esac
