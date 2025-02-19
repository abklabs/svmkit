# -*- mode: shell-script -*-
# shellcheck shell=bash

umask 077

stake-account-create () {
    solana create-stake-account "${SOLANA_CLI_TXN_FLAGS[@]}" stake_account.json "$STAKE_AMOUNT"
    solana delegate-stake "${SOLANA_CLI_TXN_FLAGS[@]}" stake_account.json vote_account.json
}

stake-account-deactivate () {
    # TODO: Add optional stake authority flag
    solana deactivate-stake "${SOLANA_CLI_TXN_FLAGS[@]}" stake_account.json
}

stake-account-withdraw () {
    # TODO: Add a withdraw-authority flag + option
    solana withdraw-stake "${SOLANA_CLI_TXN_FLAGS[@]}" stake_account.json  withdraw_account.json "$WITHDRAW_AMOUNT"
}

case "$STAKE_ACCOUNT_ACTION" in
    CREATE)
        stake-account-create
        ;;
    DEACTIVATE)
        stake-account-deactivate
        ;;
    WITHDRAW)
        stake-account-withdraw
        ;;
    *)
        log::fatal "unknown action provided '$STAKE_ACCOUNT_ACTION'!"
    ;;
esac
