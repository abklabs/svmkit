#!/usr/bin/env bash
set -euo pipefail

source ./lib.bash
source ./env
source ./steps.sh

# shellcheck disable=SC1090
steps::run "step" "$@"
