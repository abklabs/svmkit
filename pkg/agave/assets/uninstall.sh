# -*- mode: shell-script -*-
# shellcheck shell=bash

# shellcheck disable=SC1091
. ./deletion-lib.sh

VALIDATOR_SERVICE=${VALIDATOR_SERVICE}.service

step::00::wait-for-a-stable-environment() {
    cloud-init::wait-for-stable-environment
}

step::10::stop-services() {
  svmkit::sudo systemctl stop "${VALIDATOR_SERVICE}"
  svmkit::sudo systemctl disable "${VALIDATOR_SERVICE}"
}

step::80::delete-files() {
    deletion::delete
}
