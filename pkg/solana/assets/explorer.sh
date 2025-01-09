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
        svmkit::apt::get install ufw nodejs npm "svmkit-solana-explorer=$EXPLORER_VERSION"
    else
        svmkit::apt::get install ufw nodejs npm svmkit-solana-explorer
    fi
    svmkit::sudo npm install -g pnpm
}

step::004::create-sol-user() {
    create-sol-user
}

step::005::configure-firewall() {
    svmkit::sudo ufw allow "$EXPLORER_PORT/tcp"
    svmkit::sudo ufw allow 22/tcp
    svmkit::sudo ufw --force enable
    svmkit::sudo ufw reload
}

step::006::setup-explorer() {
    svmkit::sudo chown -R sol:sol /opt/svmkit-solana-explorer
    svmkit::sudo -i -u sol bash -c 'cd /opt/svmkit-solana-explorer && pnpm install'
}

step::007::setup-explorer-startup() {
    if systemctl list-unit-files "${EXPLORER_SERVICE}" >/dev/null; then
        svmkit::sudo systemctl stop "${EXPLORER_SERVICE}" || true
    fi

    cat <<EOF | svmkit::sudo tee /opt/svmkit-solana-explorer/run-explorer >/dev/null
#!/usr/bin/env bash

cd /opt/svmkit-solana-explorer
$EXPLORER_ENV exec pnpm start /opt/svmkit-solana-explorer $EXPLORER_FLAGS
EOF

    svmkit::sudo chmod 755 /opt/svmkit-solana-explorer/run-explorer
    svmkit::sudo chown sol:sol /opt/svmkit-solana-explorer/run-explorer

    cat <<EOF | svmkit::sudo tee /etc/systemd/system/"${EXPLORER_SERVICE}" >/dev/null
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
    svmkit::sudo systemctl daemon-reload
    svmkit::sudo systemctl enable "${EXPLORER_SERVICE}"
    svmkit::sudo systemctl start "${EXPLORER_SERVICE}"
}
