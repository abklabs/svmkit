# -*- mode: shell-script -*-
# shellcheck shell=bash

deletion::st_dev() {
  stat -c'%d' "$1"
}

deletion::check-create() {
  local f g
  for f in "${DELETION_CHECK_FILES[@]}"; do
    if [[ -d "$f" ]]; then
      for g in "$f"/* "$f"/.*; do
        if [[ -e "$g" ]]; then
          log::fatal "'$f' is not an empty directory - use a different deletion policy to overwrite"
        fi
      done
    elif [[ -e "$f" ]]; then
      log::fatal "'$f' already exists - use a different deletion policy to overwrite"
    fi
  done
}

deletion::delete() {
  local f g
  for f in "${DELETION_DELETE_FILES[@]}"; do
    if [[ -d "$f" && "$(deletion::st_dev "$f")" != "$(deletion::st_dev "$(dirname "$f")")" ]]; then
      for g in "$f"/* "$f"/.*; do
        if [[ -e "$g" ]]; then
          svmkit::sudo rm -rvf -- "$g"
        fi
      done
    else
      svmkit::sudo rm -rvf -- "$f"
    fi
  done
}

# vim:set ft=sh:
