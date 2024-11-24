apt::setup-abk-apt-source() {
    $APT update
    $APT install curl gnupg
    if ! grep -q "^deb .*/svmkit dev main" /etc/apt/sources.list /etc/apt/sources.list.d/*; then
        curl -s https://apt.abklabs.com/keys/abklabs-archive-dev.asc | $SUDO apt-key add -
        echo "deb https://apt.abklabs.com/svmkit dev main" | $SUDO tee /etc/apt/sources.list.d/svmkit.list >/dev/null
        $APT update
    fi
}

cloud-init::wait-for-stable-environment() {
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

create-sol-user() {
    local username

    id sol >/dev/null 2>&1 || $SUDO adduser --disabled-password --gecos "" sol
    $SUDO mkdir -p "/home/sol"
    $SUDO chown -f -R sol:sol "/home/sol"

    username=$(whoami)
    id -nGz "$username" | grep -qzxF sol || $SUDO adduser "$username" sol
}

apt::env
