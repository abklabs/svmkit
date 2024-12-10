# -*- mode: shell-script -*-
# shellcheck shell=bash

FAUCET_SERVICE=svmkit-solana-faucet.service

step::000::wait-for-a-stable-environment() {
    cloud-init::wait-for-stable-environment
}

step::001::setup-abklabs-api() {
    apt::setup-abk-apt-source
}

step::003::install-faucet() {
    $APT install ufw svmkit-solana-faucet
}

step::004::create-sol-user() {
    create-sol-user
}

step::005::configure-firewall() {
    $SUDO ufw allow $FAUCET_PORT/tcp
    $SUDO ufw allow 22/tcp
    $SUDO ufw --force enable
    $SUDO ufw reload
}

step::006::setup-faucet-startup() {
    if systemctl list-unit-files "${FAUCET_SERVICE}" >/dev/null; then
        $SUDO systemctl stop "${FAUCET_SERVICE}" || true
    fi

    cat <<EOF | $SUDO tee /home/sol/run-faucet >/dev/null
#!/usr/bin/env bash

$FAUCET_ENV exec solana-faucet $FAUCET_FLAGS
EOF

    $SUDO chmod 755 /home/sol/run-faucet
    $SUDO chown sol:sol /home/sol/run-faucet

    cat <<EOF | $SUDO tee /etc/systemd/system/"${FAUCET_SERVICE}" >/dev/null
[Unit]
Description=SVMkit Solana Faucet

[Service]
Type=exec
User=sol
Group=sol
ExecStart=/home/sol/run-faucet

[Install]
WantedBy=default.target
EOF
    $SUDO systemctl daemon-reload
    $SUDO systemctl enable "${FAUCET_SERVICE}"
    $SUDO systemctl start "${FAUCET_SERVICE}"
}
