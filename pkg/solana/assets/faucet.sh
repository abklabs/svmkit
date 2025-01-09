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
    if [[ -v FAUCET_VERSION ]]; then
        svmkit::apt::get install ufw "svmkit-solana-faucet=$FAUCET_VERSION"
    else
        svmkit::apt::get install ufw svmkit-solana-faucet
    fi
}

step::004::create-sol-user() {
    create-sol-user
}

step::005::configure-firewall() {
    svmkit::sudo ufw allow "$FAUCET_PORT/tcp"
    svmkit::sudo ufw allow 22/tcp
    svmkit::sudo ufw --force enable
    svmkit::sudo ufw reload
}

step::006::copy-faucet-keys() {
    svmkit::sudo cp faucet-keypair.json /home/sol
    svmkit::sudo chown sol:sol /home/sol/faucet-keypair.json
}

step::007::setup-faucet-startup() {
    if systemctl list-unit-files "${FAUCET_SERVICE}" >/dev/null; then
        svmkit::sudo systemctl stop "${FAUCET_SERVICE}" || true
    fi

    cat <<EOF | svmkit::sudo tee /home/sol/run-faucet >/dev/null
#!/usr/bin/env bash

$FAUCET_ENV exec solana-faucet $FAUCET_FLAGS
EOF

    svmkit::sudo chmod 755 /home/sol/run-faucet
    svmkit::sudo chown sol:sol /home/sol/run-faucet

    cat <<EOF | svmkit::sudo tee /etc/systemd/system/"${FAUCET_SERVICE}" >/dev/null
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
    svmkit::sudo systemctl daemon-reload
    svmkit::sudo systemctl enable "${FAUCET_SERVICE}"
    svmkit::sudo systemctl start "${FAUCET_SERVICE}"
}
