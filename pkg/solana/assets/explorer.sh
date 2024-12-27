# -*- mode: shell-script -*-
# shellcheck shell=bash

EXPLORER_SERVICE=svmkit-solana-explorer.service

step::000::wait-for-a-stable-environment() {
    cloud-init::wait-for-stable-environment
}

step::001::setup-abklabs-api() {
    apt::setup-abk-apt-source
}

step::003::install-explorer() {
    if [[ -v EXPLORER_VERSION ]]; then
        $APT install ufw nodejs npm "svmkit-solana-explorer=$EXPLORER_VERSION"
    else
        $APT install ufw nodejs npm svmkit-solana-explorer
    fi
    $SUDO npm install -g pnpm
}

step::004::create-sol-user() {
    create-sol-user
}

step::005::configure-firewall() {
    $SUDO ufw allow "$EXPLORER_PORT/tcp"
    $SUDO ufw allow 22/tcp
    $SUDO ufw --force enable
    $SUDO ufw reload
}

step::006::setup-explorer() {
    $SUDO chown -R sol:sol /opt/svmkit-solana-explorer
    $SUDO -i -u sol bash -c 'cd /opt/svmkit-solana-explorer && pnpm install'
}

step::007::setup-explorer-startup() {
    if systemctl list-unit-files "${EXPLORER_SERVICE}" >/dev/null; then
        $SUDO systemctl stop "${EXPLORER_SERVICE}" || true
    fi

    cat <<EOF | $SUDO tee /opt/svmkit-solana-explorer/run-explorer >/dev/null
#!/usr/bin/env bash

cd /opt/svmkit-solana-explorer
$EXPLORER_ENV exec pnpm start /opt/svmkit-solana-explorer $EXPLORER_FLAGS
EOF

    $SUDO chmod 755 /opt/svmkit-solana-explorer/run-explorer
    $SUDO chown sol:sol /opt/svmkit-solana-explorer/run-explorer

    cat <<EOF | $SUDO tee /etc/systemd/system/"${EXPLORER_SERVICE}" >/dev/null
[Unit]
Description=SVMkit Solana Explorer

[Service]
Type=exec
User=sol
Group=sol
ExecStart=/opt/svmkit-solana-explorer/run-explorer

[Install]
WantedBy=default.target
EOF
    $SUDO systemctl daemon-reload
    $SUDO systemctl enable "${EXPLORER_SERVICE}"
    $SUDO systemctl start "${EXPLORER_SERVICE}"
}
