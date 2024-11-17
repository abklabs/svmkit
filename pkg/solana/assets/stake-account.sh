# -*- mode: shell-script -*-
# shellcheck shell=bash

umask 077

stake-account-create () {
    local stake_account_keypair vote_account_keypair

    stake_account_keypair=$(temp::file)
    # shellcheck disable=SC2153
    echo "$STAKE_ACCOUNT_KEYPAIR" > "$stake_account_keypair"

    vote_account_keypair=$(temp::file)
    # shellcheck disable=SC2153
    echo "$VOTE_ACCOUNT_KEYPAIR" > "$vote_account_keypair"

    solana create-stake-account "$stake_account_keypair" "$STAKE_AMOUNT"
    solana delegate-stake "$stake_account_keypair" "$vote_account_keypair"
}

case "$STAKE_ACCOUNT_ACTION" in
    CREATE)
	stake-account-create
	;;
    *)
	log::fatal "unknown action provided '$STAKE_ACCOUNT_ACTION'!"
    ;;
esac
