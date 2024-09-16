#!/usr/bin/env bash
set -euo pipefail

source ./lib.bash
source ./env
source ./steps.sh

step_prefix="step"

run() {
  local start name

  start=""

  if [[ $# -gt 0 ]]; then
    start="${step_prefix}::$1"
    shift
    log::warn "starting steps with $start..."
  fi

  while read -r name; do
    if [[ $name > $start || $name = "$start" ]]; then
      log::info "====> running $name..."
      $name
    fi
  done < <(declare -F | grep "${step_prefix}" | awk '{ print $3; }')
}

# shellcheck disable=SC1090
run "$@"
