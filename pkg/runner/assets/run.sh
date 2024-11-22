#!/usr/bin/env ./opsh

source ./lib.bash
source ./env
source ./steps.sh

# shellcheck disable=SC1090
steps::run "step" "$@"
