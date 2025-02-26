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

stake-account-update () {
    local stake_auth_args=()
    if [[ -v STAKE_AUTHORITY ]]; then
      stake_auth_args+=(--stake-authority stake_authority.json)
    fi

    local withdraw_auth_args=()
    if [[ -v WITHDRAW_AUTHORITY ]]; then
      withdraw_auth_args+=(--withdraw-authority withdraw_authority.json)
    fi

    # Handle deactivation or delegation (mutually exclusive)
    if [[ -v STAKE_ACCOUNT_DEACTIVATE ]] && [[ "$STAKE_ACCOUNT_DEACTIVATE" == "true" ]]; then
        solana deactivate-stake "${SOLANA_CLI_TXN_FLAGS[@]}" stake_account.json "${stake_auth_args[@]}"
    elif [[ -v STAKE_ACCOUNT_DELEGATE ]] && [[ "$STAKE_ACCOUNT_DELEGATE" == "true" ]]; then
        solana delegate-stake "${SOLANA_CLI_TXN_FLAGS[@]}" stake_account.json new_vote_account.json "${stake_auth_args[@]}"
    fi

    # Handle authority changes if requested
    if [[ -v STAKE_ACCOUNT_AUTHORITY ]] && [[ "$STAKE_ACCOUNT_AUTHORITY" == "true" ]]; then
        if [[ -v STAKE_AUTHORITY_UPDATE ]] && [[ "$STAKE_AUTHORITY_UPDATE" == "true" ]] && [[ -f new_stake_authority.json ]]; then
            solana stake-authorize "${SOLANA_CLI_TXN_FLAGS[@]}" \
                stake_account.json \
                --stake-authority stake_authority.json \
                --new-stake-authority new_stake_authority.json
        fi
        
        if [[ -v WITHDRAW_AUTHORITY_UPDATE ]] && [[ "$WITHDRAW_AUTHORITY_UPDATE" == "true" ]] && [[ -f new_withdraw_authority.json ]]; then
            solana stake-authorize "${SOLANA_CLI_TXN_FLAGS[@]}" \
                stake_account.json \
                --withdraw-authority withdraw_authority.json \
                --new-withdraw-authority new_withdraw_authority.json
        fi
    fi

    # Handle lockup changes if requested
    if [[ -v STAKE_ACCOUNT_LOCKUP ]] && [[ "$STAKE_ACCOUNT_LOCKUP" == "true" ]] && [[ -v CUSTODIAN_PUBKEY ]] && [[ -v EPOCH_AVAILABLE ]]; then
        solana stake-set-lockup "${SOLANA_CLI_TXN_FLAGS[@]}" stake_account.json \
            --custodian "$CUSTODIAN_PUBKEY" \
            --epoch "$EPOCH_AVAILABLE" \
            "${withdraw_auth_args[@]}"
    fi
}

case "$STAKE_ACCOUNT_ACTION" in
    CREATE)
	stake-account-create
	;;
    DELETE)
	stake-account-delete
	;;
    UPDATE)
	stake-account-update
	;;
    *)
	log::fatal "unknown action provided '$STAKE_ACCOUNT_ACTION'!"
    ;;
esac
