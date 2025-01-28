svmkit::sudo() {
    sudo "$@"
}

svmkit::apt::get() {
    svmkit::sudo DEBIAN_FRONTEND=noninteractive apt-get -qy -o APT::Lock::Timeout="${APT_LOCK_TIMEOUT}" -o DPkg::Lock::Timeout="${APT_LOCK_TIMEOUT}" "$@"
}

wait-for-update-lock() {
    # If apt-get update gets called in parallel too closely, it can cause a
    # lock contention that the Lock::Timeout can't catch. This is a workaround
    # to wait for the lock to be released.
    local sleep_time=5
    local tries=0
    local max_tries=$(((APT_LOCK_TIMEOUT + sleep_time - 1) / sleep_time))
    while pgrep -x apt-get >/dev/null; do
        log::warn "Another instance of apt-get update is running. Waiting..."
        sleep $sleep_time
        tries=$((tries + 1))
        if [[ $tries -gt $max_tries ]]; then
            log::warn "Waited $((tries * sleep_time)) seconds total, but apt-get is still running. Giving up."
            break
        fi
    done
}

svmkit::apt::update() {
    wait-for-update-lock
    svmkit::apt::get update
}

apt::setup-abk-apt-source() {
    svmkit::apt::update

    svmkit::apt::get install curl gnupg
    if ! grep -q "^deb .*/svmkit dev main" /etc/apt/sources.list /etc/apt/sources.list.d/*; then
        curl -s https://apt.abklabs.com/keys/abklabs-archive-dev.asc | svmkit::sudo apt-key add -
        echo "deb https://apt.abklabs.com/svmkit dev main" | svmkit::sudo tee /etc/apt/sources.list.d/svmkit.list >/dev/null

        svmkit::apt::update
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
