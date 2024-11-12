# -*- mode: shell-script -*-
# shellcheck shell=bash

umask 077

vote-account-create () {
    local identity_keypair vote_account_keypair auth_withdrawer_keypair args

    identity_keypair=$(temp::file)
    echo "$IDENTITY_KEYPAIR" > "$identity_keypair"

    vote_account_keypair=$(temp::file)
    echo "$VOTE_ACCOUNT_KEYPAIR" > "$vote_account_keypair"

    auth_withdrawer_keypair=$(temp::file)
    echo "$AUTH_WITHDRAWER_KEYPAIR" > "$auth_withdrawer_keypair"

    args=()

    if [[ -v AUTH_VOTER_PUBKEY ]]; then
	args+=(--authorized-voter "$AUTH_VOTER_PUBKEY")
    fi

    solana create-vote-account "$vote_account_keypair" "$identity_keypair" "$auth_withdrawer_keypair"
}

vote-account-delete () {
    local identity_keypair vote_account_keypair auth_withdrawer_keypair

    # If they haven't provided a close-recipient public key, then
    # don't bother closing down the account.  The logic being either:
    #
    # 1) That the vote account is either going to get thrown away
    # (e.g. if the entire cluster is being torn down, in the case of
    # an SPE).
    #
    # 2) The user doesn't want a closed account because they'll be
    # specifying the close recipient manually.

    [[ -v CLOSE_RECIPIENT_PUBKEY ]] || return 0

    vote_account_keypair=$(temp::file)
    echo "$VOTE_ACCOUNT_KEYPAIR" > "$vote_account_keypair"

    auth_withdrawer_keypair=$(temp::file)
    echo "$AUTH_WITHDRAWER_KEYPAIR" > "$auth_withdrawer_keypair"

    solana close-vote-account --authorized-withdrawer "$auth_withdrawer_keypair" "$vote_account_keypair" "$CLOSE_RECIPIENT_PUBKEY"
}

case "$VOTE_ACCOUNT_ACTION" in
    CREATE)
	vote-account-create
	;;
    DELETE)
	vote-account-delete
	;;
    *)
	log::fatal "unknown action provided '$VOTE_ACCOUNT_ACTION'!"
    ;;
esac
