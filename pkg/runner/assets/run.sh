#!/usr/bin/env ./opsh

PATH="$PWD:$PATH"

source ./lib.bash
source ./env
source ./steps.sh

# shellcheck disable=SC1090
steps::run "step" "$@"
