# -*- mode: shell-script -*-
# shellcheck shell=bash

step::00::update-run-validator() {
  cat <<EOF | $SUDO tee /home/sol/run-validator >/dev/null
#!/usr/bin/env bash
exec agave-validator $VALIDATOR_FLAGS
EOF

  $SUDO systemctl restart svmkit-agave-validator.service
}
