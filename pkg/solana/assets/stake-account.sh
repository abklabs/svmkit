# -*- mode: shell-script -*-
# shellcheck shell=bash

umask 077

stake-account-create () {
    local create_args=()
    local delegate_args=()

    if [[ -v WITHDRAW_AUTHORITY ]]; then
      create_args+=(--withdraw-authority withdraw_authority.json)
    fi

    if [[ -v STAKE_AUTHORITY ]]; then
      create_args+=(--stake-authority stake_authority.json)
      delegate_args+=(--stake-authority stake_authority.json)
    fi

    solana create-stake-account "${SOLANA_CLI_TXN_FLAGS[@]}" stake_account.json "$STAKE_AMOUNT" "${create_args[@]}"
    solana delegate-stake "${SOLANA_CLI_TXN_FLAGS[@]}" stake_account.json vote_account.json "${delegate_args[@]}"
}

stake-account-delete () {
    [[ -v FORCE_DELETE ]] && [[ "$FORCE_DELETE" == "true" ]] || return 0

    local args=()
    if [[ -v WITHDRAW_AUTHORITY ]]; then
      args+=(--withdraw-authority withdraw_authority.json)
    fi

    solana withdraw "${SOLANA_CLI_TXN_FLAGS[@]}" stake_account.json $WITHDRAW_PUBKEY "$STAKE_AMOUNT" "${args[@]}"
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
