# -*- mode: shell-script -*-
# shellcheck shell=bash

umask 077

stake-account-status () {
    solana stake-account --output json-compact "$STAKE_ACCOUNT_ADDRESS"
}

stake-account-create () {
    local args=()

    if [[ -n "${STAKE_AUTHORITY_ADDRESS}" ]]; then
      args+=(--stake-authority "$STAKE_AUTHORITY_ADDRESS")
    fi

    if [[ -n "${WITHDRAW_AUTHORITY_ADDRESS}" ]]; then
      args+=(--withdraw-authority "$WITHDRAW_AUTHORITY_ADDRESS")
    fi

    if [[ -v CUSTODIAN_PUBKEY ]] && [[ -v EPOCH_AVAILABLE ]] && [[ -n "$CUSTODIAN_PUBKEY" ]] && [[ -n "$EPOCH_AVAILABLE" ]]; then
      args+=(--lockup-epoch "$EPOCH_AVAILABLE" --custodian "$CUSTODIAN_PUBKEY")
    fi
    solana create-stake-account "${SOLANA_CLI_TXN_FLAGS[@]}" stake_account.json "$STAKE_AMOUNT" "${args[@]}"
}

stake-account-deactivate () {
  if [[ -f stake_authority.json ]]; then
    solana deactivate-stake "${SOLANA_CLI_TXN_FLAGS[@]}" "$STAKE_ACCOUNT_ADDRESS" --stake-authority stake_authority.json
  else
    solana deactivate-stake "${SOLANA_CLI_TXN_FLAGS[@]}" "$STAKE_ACCOUNT_ADDRESS"
  fi
}

stake-account-delegate () {
    local delegate_args=()
    if [[ -f stake_authority.json ]]; then
      delegate_args+=(--stake-authority stake_authority.json)
    fi

    solana delegate-stake "${SOLANA_CLI_TXN_FLAGS[@]}" "$STAKE_ACCOUNT_ADDRESS" vote_account.json "${delegate_args[@]}"
}

stake-account-authorize() {
    local authorize_args=()

    if [[ $AUTH_TYPE == "STAKER" ]]; then
        authorize_args+=(--new-stake-authority new_address.json)
        if [[ -f old_address.json ]]; then
          authorize_args+=(--stake-authority old_address.json)
        fi

    elif [[ $AUTH_TYPE == "WITHDRAWER" ]]; then
        authorize_args+=(--new-withdraw-authority new_address.json)
        if [[ -f old_address.json ]]; then
          authorize_args+=(--withdraw-authority old_address.json)
        fi
        if [[ -f lockup_keypair.json ]]; then
          authorize_args+=(--lockup-authority lockup_keypair.json)
        fi
    else
        echo "Error: Invalid authorization type provided"
        return 1
    fi

    solana stake-authorize "${SOLANA_CLI_TXN_FLAGS[@]}" "$STAKE_ACCOUNT_ADDRESS" "${authorize_args[@]}"
}

stake-account-withdraw() {
    local args=()
    if [[ -f withdraw_authority.json ]]; then
      args+=(--withdraw-authority withdraw_authority.json)
    fi

    if [[ -f lockup_authority.json ]]; then
      args+=(--lockup-authority lockup_authority.json)
    fi

    solana withdraw-stake "${SOLANA_CLI_TXN_FLAGS[@]}" "$STAKE_ACCOUNT_ADDRESS" "$WITHDRAW_ADDRESS" "$STAKE_AMOUNT" "${args[@]}"
}

stake-account-set-lockup() {
  log::fatal "Not implemented"
}

case "$STAKE_ACCOUNT_ACTION" in
    READ)
	stake-account-status
	;;
    CREATE)
	stake-account-create
	;;
    DELEGATE)
	stake-account-delegate
	;;
    DEACTIVATE)
	stake-account-deactivate
	;;
    AUTHORIZE)
	stake-account-authorize
	;;
    WITHDRAW)
	stake-account-withdraw
	;;
    LOCKUP)
	stake-account-set-lockup
	;;
    *)
	log::fatal "unknown action provided '$STAKE_ACCOUNT_ACTION'!"
    ;;
esac
