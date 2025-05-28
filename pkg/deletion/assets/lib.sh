# -*- mode: shell-script -*-
# shellcheck shell=bash

deletion::check-create() {
  local f
  for f in "${DELETION_CHECK_FILES[@]}"; do
    if [[ -e "$f" ]]; then
      echo "'$f' already exists - use a different deletion policy to overwrite"
      return 255
    fi
  done
}

deletion::delete() {
  local f
  for f in "${DELETION_DELETE_FILES[@]}"; do
    rm -rvf -- "$f"
  done
}

# vim:set ft=sh:
