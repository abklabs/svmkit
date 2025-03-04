# -*- mode: shell-script -*-
# shellcheck shell=bash

umask 077

stake-account-status () {
    solana stake-account --output json-compact stake_account.json
}

stake-account-create () {
    local create_args=()
    local delegate_args=()

    # Check for authority files directly
    if [[ -f withdraw_authority.json ]]; then
      create_args+=(--withdraw-authority withdraw_authority.json)
    fi

    if [[ -f stake_authority.json ]]; then
      create_args+=(--stake-authority stake_authority.json)
      delegate_args+=(--stake-authority stake_authority.json)
    fi

    # Lockup args still need environment variables
    if [[ -v CUSTODIAN_PUBKEY ]] && [[ -v EPOCH_AVAILABLE ]] && [[ -n "$CUSTODIAN_PUBKEY" ]] && [[ -n "$EPOCH_AVAILABLE" ]]; then
      create_args+=(--lockup-epoch "$EPOCH_AVAILABLE" --custodian "$CUSTODIAN_PUBKEY")
    fi

    solana create-stake-account "${SOLANA_CLI_TXN_FLAGS[@]}" stake_account.json "$STAKE_AMOUNT" "${create_args[@]}"
    
    # Only delegate if vote account exists
    if [[ -f vote_account.json ]]; then
        solana delegate-stake "${SOLANA_CLI_TXN_FLAGS[@]}" stake_account.json vote_account.json "${delegate_args[@]}"
    fi
}


stake-account-delete () {
    # Return early if force delete is not set to true
    [[ "$FORCE_DELETE" == "true" ]] || return 0
    
    # Check if we have a withdraw address (ensuring it exists first)
    if [[ ! -v WITHDRAW_ADDRESS ]] || [[ -z "$WITHDRAW_ADDRESS" ]]; then
        echo "Error: No withdraw address provided for stake account deletion"
        return 1
    fi
    
    local args=()
    if [[ -f withdraw_authority.json ]]; then
      args+=(--withdraw-authority withdraw_authority.json)
    fi

    solana withdraw "${SOLANA_CLI_TXN_FLAGS[@]}" stake_account.json "$WITHDRAW_ADDRESS" "$STAKE_AMOUNT" "${args[@]}"
}

stake-account-update () {
    # Build authority argument arrays based on file existence
    local stake_auth_args=()
    if [[ -f stake_authority.json ]]; then
      stake_auth_args+=(--stake-authority stake_authority.json)
    fi

    local withdraw_auth_args=()
    if [[ -f withdraw_authority.json ]]; then
      withdraw_auth_args+=(--withdraw-authority withdraw_authority.json)
    fi

    # Handle deactivation (required explicit operation flag as it doesn't use any unique files)
    if [[ "$OPERATION" == "DEACTIVATE" ]]; then
        solana deactivate-stake "${SOLANA_CLI_TXN_FLAGS[@]}" stake_account.json "${stake_auth_args[@]}"
    fi

    # Handle delegation (check for new vote account file)
    if [[ -f new_vote_account.json ]]; then
        solana delegate-stake "${SOLANA_CLI_TXN_FLAGS[@]}" stake_account.json new_vote_account.json "${stake_auth_args[@]}"
    fi

    # Handle stake authority update
    if [[ "$UPDATE_STAKE_AUTHORITY" == "true" ]] && [[ -f stake_authority.json ]] && [[ -f new_stake_authority.json ]]; then
        solana stake-authorize "${SOLANA_CLI_TXN_FLAGS[@]}" \
            stake_account.json \
            --stake-authority stake_authority.json \
            --new-stake-authority new_stake_authority.json
    fi

    # Handle withdraw authority update
    if [[ "$UPDATE_WITHDRAW_AUTHORITY" == "true" ]] && [[ -f withdraw_authority.json ]] && [[ -f new_withdraw_authority.json ]]; then
        solana stake-authorize "${SOLANA_CLI_TXN_FLAGS[@]}" \
            stake_account.json \
            --withdraw-authority withdraw_authority.json \
            --new-withdraw-authority new_withdraw_authority.json
    fi

    # Handle lockup changes (still needs explicit params)
    if [[ -v CUSTODIAN_PUBKEY ]] && [[ -v EPOCH_AVAILABLE ]] && [[ -n "$CUSTODIAN_PUBKEY" ]] && [[ -n "$EPOCH_AVAILABLE" ]]; then
        solana stake-set-lockup "${SOLANA_CLI_TXN_FLAGS[@]}" stake_account.json \
            --custodian "$CUSTODIAN_PUBKEY" \
            --epoch "$EPOCH_AVAILABLE" \
            "${withdraw_auth_args[@]}"
    fi
}

case "$STAKE_ACCOUNT_ACTION" in
    READ)
	stake-account-status
	;;
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
