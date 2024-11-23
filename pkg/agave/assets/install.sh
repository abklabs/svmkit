# -*- mode: shell-script -*-
# shellcheck shell=bash

: "${RPC_SERVICE_TIMEOUT:=60}"

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
    local ret

    if command -v cloud-init >/dev/null 2>&1; then
        if systemctl is-active --quiet cloud-init.service; then
            ret=0
            cloud-init status --wait || ret=$?

            case "$ret" in
                0)
                    log::info "cloud-init has finished, continuing on"
                    ;;
                2)
                    log::warn "cloud-init had a recoverable error; we're continuing anyway"
                    ;;
                *)
                    log::error "cloud-init status exited with status $ret; continuing but you should investigate"
                    ;;
            esac
        else
            log::warn "cloud-init.service in a failed state; not waiting for completion"
        fi
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
    local username

    id sol >/dev/null 2>&1 || $SUDO adduser --disabled-password --gecos "" sol
    $SUDO mkdir -p "/home/sol"
    $SUDO chown -f -R sol:sol "/home/sol"

    username=$(whoami)
    id -nGz "$username" | grep -qzxF sol || $SUDO adduser "$username" sol
}

step::30::copy-validator-keys() {
    $SUDO cp validator-keypair.json vote-account-keypair.json /home/sol
    $SUDO chown sol:sol /home/sol/{validator-keypair,vote-account-keypair}.json
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

    cat <<EOF | $SUDO tee /home/sol/check-validator >/dev/null
#!/usr/bin/env bash
set -euo pipefail

FULL_RPC=${FULL_RPC:=false}
RPC_BIND_ADDRESS=$RPC_BIND_ADDRESS
RPC_PORT=$RPC_PORT
RPC_SERVICE_TIMEOUT=$RPC_SERVICE_TIMEOUT

\$FULL_RPC || exit 0

for i in \$(seq 1 \$RPC_SERVICE_TIMEOUT) ; do
    if solana slot --url http://\$RPC_BIND_ADDRESS:\$RPC_PORT &> /dev/null ; then
        exit 0
    fi
    sleep 1
done

echo "timed out waiting for validator to bring RPC online!" 1>&2
exit 1
EOF

    for i in run-validator check-validator ; do
	$SUDO chmod 755 /home/sol/$i
	$SUDO chown sol:sol /home/sol/$i
    done

    cat <<EOF | $SUDO tee /etc/systemd/system/"${VALIDATOR_SERVICE}" >/dev/null
[Unit]
Description=SVMkit $VALIDATOR_VARIANT validator

[Service]
Type=exec
User=sol
Group=sol
ExecStart=/home/sol/run-validator
ExecStartPost=/home/sol/check-validator
LimitNOFILE=1000000

[Install]
WantedBy=default.target
EOF
    $SUDO systemctl daemon-reload
    $SUDO systemctl enable "${VALIDATOR_SERVICE}"
    $SUDO systemctl start "${VALIDATOR_SERVICE}"
}

step::90::setup-validator-info() {
    local args

    [[ -v VALIDATOR_INFO_NAME ]] || return 0

    if [[ -v VALIDATOR_INFO_WEBSITE ]] ; then
	args+=(--website "$VALIDATOR_INFO_WEBSITE")
    fi

    if [[ -v VALIDATOR_INFO_ICON_URL ]] ; then
	args+=(--icon-url "$VALIDATOR_INFO_ICON_URL")
    fi

    if [[ -v VALIDATOR_INFO_DETAILS ]] ; then
	args+=(--details "$VALIDATOR_INFO_DETAILS")
    fi

    $SUDO -u sol -i solana validator-info publish "${args[@]}" "$VALIDATOR_INFO_NAME"
}
