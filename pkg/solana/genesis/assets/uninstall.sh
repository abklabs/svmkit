# -*- mode: shell-script -*-
# shellcheck shell=bash

# shellcheck disable=SC1091
. ./deletion-lib.sh

step::080::delete-files() {
    deletion::delete
}
