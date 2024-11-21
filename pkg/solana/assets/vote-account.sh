# -*- mode: shell-script -*-
# shellcheck shell=bash

umask 077

vote-account-create () {
    local args=()

    if [[ -v AUTH_VOTER_PUBKEY ]]; then
	args+=(--authorized-voter "$AUTH_VOTER_PUBKEY")
    fi

    solana create-vote-account vote_account.json identity.json auth_withdrawer.json
}

vote-account-delete () {
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

    solana close-vote-account --authorized-withdrawer auth_withdrawer.json vote_account.json "$CLOSE_RECIPIENT_PUBKEY"
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
