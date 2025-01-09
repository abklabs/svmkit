svmkit::sudo () {
    sudo "$@"
}

svmkit::apt::get() {
    svmkit::sudo DEBIAN_FRONTEND=noninteractive apt-get -qy "$@"
}

apt::setup-abk-apt-source() {
    svmkit::apt::get update
    svmkit::apt::get install curl gnupg
    if ! grep -q "^deb .*/svmkit dev main" /etc/apt/sources.list /etc/apt/sources.list.d/*; then
        curl -s https://apt.abklabs.com/keys/abklabs-archive-dev.asc | svmkit::sudo apt-key add -
        echo "deb https://apt.abklabs.com/svmkit dev main" | svmkit::sudo  tee /etc/apt/sources.list.d/svmkit.list >/dev/null
        svmkit::apt::get update
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

    id sol >/dev/null 2>&1 || svmkit::sudo adduser --disabled-password --gecos "" sol
    svmkit::sudo mkdir -p "/home/sol"
    svmkit::sudo chown -f -R sol:sol "/home/sol"

    username=$(whoami)
    id -nGz "$username" | grep -qzxF sol || svmkit::sudo adduser "$username" sol
}
