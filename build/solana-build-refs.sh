#!/usr/bin/env bash

set -euo pipefail

log::generic() {
    local level
    level=$1
    shift

    printf "%s\t%s\n" "$level" "$*"
}

log::info() {
    log::generic INFO "$@"
}

log::fatal() {
    log::generic FATAL "$@"
    exit 1
}

lookup-remote-tag() {
    local remote tag tagfile tagcount
    remote=$1
    shift
    tag=$1
    shift

    tagfile=$(mktemp)

    git ls-remote --tags "$remote" "$tag" >"$tagfile"

    tagcount=$(wc -l <"$tagfile")

    if [[ $tagcount -lt 1 ]]; then
        log::fatal "no tags found on $remote for $tag!"
    fi

    if [[ $tagcount -gt 1 ]]; then
        log::fatal "found more than one tag matching $tag on $remote.  cowardly giving up!"
    fi

    awk '{ print $1;}' <"$tagfile"
    rm "$tagfile"
}

fetch-remote() {
    local remote
    remote=$1
    shift
    log::info "git fetching remote $remote..."
    git fetch "$remote"
}

lookup-grpc-ref() {
    local agave_tag=$1
    local pattern="*+solana.${agave_tag#v}"
    readarray -t matches < <(
        git ls-remote --tags yellowstone-grpc "$pattern" |
        awk '{print $2}' | sed 's|refs/tags/||'
    )
    if [[ ${#matches[@]} -eq 0 ]]; then
        log::info "no grpc tag for $pattern – skipping"
        return 1
    elif [[ ${#matches[@]} -gt 1 ]]; then
        log::fatal "multiple grpc tags for $pattern"
    fi
    echo "${matches[0]}"
}

geyser-interface-version() {
    local commit=$1
    git checkout -f "$commit"
    cargo metadata --format-version 1 \
        --manifest-path validator/Cargo.toml |
    jq -r '.packages[] | select(.name=="agave-geyser-plugin-interface") | .version'
}

for remote in solana-labs anza-xyz PowerLedger jito-foundation pyth-network mantis xen tachyon yellowstone-grpc; do
    fetch-remote $remote
done

build-ref svmkit-solana solana-labs/master solana-validator

for tag in v2.1.13 v2.1.14 v2.1.15 v2.1.16 v2.1.21 v2.2.0 v2.2.1; do
    commit=$(lookup-remote-tag anza-xyz "$tag")
    agave_ver=$(geyser-interface-version "$commit")
    grpc_ref=$(lookup-grpc-ref "$tag")
    if [[ -z $grpc_ref ]]; then
        log::info "no grpc tag for $tag – skipping"
        continue
    else
        log::info "found grpc tag $grpc_ref for $tag"
    fi
done
