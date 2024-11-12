# -*- mode: shell-script -*-
# shellcheck shell=bash

VALIDATOR_PACKAGE=svmkit-${VALIDATOR_VARIANT}-validator
VALIDATOR_SERVICE=${VALIDATOR_PACKAGE}.service

case $VALIDATOR_VARIANT in
    agave)
	VALIDATOR_PROCESS=agave-validator
	;;
    jito)
	VALIDATOR_PROCESS=agave-validator
	;;
    mantis)
	VALIDATOR_PROCESS=solana-validator
	;;
    powerledger)
	VALIDATOR_PROCESS=solana-validator
	;;
    pyth)
	VALIDATOR_PROCESS=solana-validator
	;;
    solana)
	VALIDATOR_PROCESS=solana-validator
	;;
    *)
	log::fatal "unknown validator variant '$VALIDATOR_VARIANT'!"
	;;
esac

step::00::wait-for-a-stable-environment() {
    if command -v cloud-init > /dev/null 2>&1 ; then
	cloud-init status --wait
    fi
}

step::05::setup-abklabs-apt() {
    apt::abk
}

step::10::install-base-software() {
    $SUDO apt-get update
    $APT install logrotate ufw
}

step::20::create-sol-user() {
    id sol >/dev/null 2>&1 || $SUDO adduser --disabled-password --gecos "" sol
    $SUDO mkdir -p "/home/sol"
    $SUDO chown -f -R sol:sol "/home/sol"
}

step::30::copy-validator-keys() {
    keypairs::write "sol" "${IDENTITY_KEYPAIR}:validator-keypair" "${VOTE_ACCOUNT_KEYPAIR}:vote-account-keypair"
}

step::40::configure-sysctl() {
    cat <<EOF | $SUDO tee /etc/sysctl.d/21-solana-validator.conf >/dev/null
# Increase UDP buffer sizes
net.core.rmem_default = 134217728
net.core.rmem_max = 134217728
net.core.wmem_default = 134217728
net.core.wmem_max = 134217728
# Increase memory mapped files limit
vm.max_map_count = 1000000
# Increase number of allowed open file descriptors
fs.nr_open = 1000000
vm.swappiness=1
EOF

    $SUDO sysctl -p /etc/sysctl.d/21-solana-validator.conf
}

step::50::configure-firewall() {
    $SUDO ufw allow 53
    $SUDO ufw allow ssh
    $SUDO ufw allow 8000:8020/tcp
    $SUDO ufw allow 8000:8020/udp
    # TODO: Only open for RPC nodes
    $SUDO ufw allow 8899/tcp
    $SUDO ufw allow 8899/udp
    $SUDO ufw allow 8900/tcp

    $SUDO ufw --force enable
}

step::60::setup-logrotate() {
    cat <<EOF | $SUDO tee /etc/logrotate.d/solana >/dev/null
/home/sol/solana-validator.log {
su sol sol
daily
rotate 1
missingok
postrotate
    systemctl kill -s USR1 sol.service
endscript
}
EOF

    $SUDO systemctl restart logrotate
}

step::70::install-validator() {
    if [[ -v VALIDATOR_VERSION ]]; then
        $APT --allow-downgrades install "${VALIDATOR_PACKAGE}=$VALIDATOR_VERSION" "svmkit-solana-cli=$VALIDATOR_VERSION"
    else
        $APT --allow-downgrades install "${VALIDATOR_PACKAGE}" "svmkit-solana-cli"
    fi
}

step::75::setup-solana-cli() {
    [[ -v SOLANA_CLI_CONFIG_FLAGS ]] || return 0

    # First setup the login user.
    solana config set $SOLANA_CLI_CONFIG_FLAGS

    # Setup the sol user.
    $SUDO -u sol -i solana config set $SOLANA_CLI_CONFIG_FLAGS
}

step::80::setup-validator-startup() {
    if systemctl list-unit-files "${VALIDATOR_SERVICE}" >/dev/null; then
        $SUDO systemctl stop "${VALIDATOR_SERVICE}" || true
    fi

    cat <<EOF | $SUDO tee /home/sol/run-validator >/dev/null
#!/usr/bin/env bash

$VALIDATOR_ENV exec $VALIDATOR_PROCESS $VALIDATOR_FLAGS
EOF

    $SUDO chmod 755 /home/sol/run-validator
    $SUDO chown sol:sol /home/sol/run-validator

    cat <<EOF | $SUDO tee /etc/systemd/system/"${VALIDATOR_SERVICE}" >/dev/null
[Unit]
Description=SVMkit $VALIDATOR_VARIANT validator

[Service]
User=sol
Group=sol
ExecStart=/home/sol/run-validator
LimitNOFILE=1000000

[Install]
WantedBy=default.target
EOF
    $SUDO systemctl daemon-reload
    $SUDO systemctl enable "${VALIDATOR_SERVICE}"
    $SUDO systemctl start "${VALIDATOR_SERVICE}"
}
