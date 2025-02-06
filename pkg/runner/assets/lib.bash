svmkit::sudo() {
    sudo "$@"
}

APT_LOCKFILE="/var/lib/dpkg/apt-svmkit.lock"

svmkit::sudo touch "$APT_LOCKFILE"
svmkit::sudo chown "$(id -u):$(id -g)" "$APT_LOCKFILE"
svmkit::sudo chmod 600 "$APT_LOCKFILE"

svmkit::flock::start() {
    if [[ -n "${lock_fd:-}" ]]; then
        log::error "Cannot reacquire svmkit lock: lock_fd=$lock_fd is already set."
        return 1
    fi

    log::info "Acquiring svmkit lock..."

    # shellcheck disable=SC2093,SC1083
    exec {lock_fd}>>"${APT_LOCKFILE}"
    flock -x -w "${APT_LOCK_TIMEOUT}" "${lock_fd}" || {
        log::error "Could not acquire svmkit lock within ${APT_LOCK_TIMEOUT}"
        exit 1
    }
}

svmkit::flock::end() {
    if [[ -z "${lock_fd:-}" ]]; then
        log::warn "Svmkit lock already released, nothing to do."
        return 0
    fi

    log::info "Releasing svmkit lock..."

    flock -u "${lock_fd}"
    exec {lock_fd}>&-
    unset lock_fd
}

svmkit::flock::cleanup() {
    [[ -v lock_fd ]] || return 0

    svmkit::flock::end
}

exit::trigger svmkit::flock::cleanup

svmkit::flock::run() {
    if [[ -n "${lock_fd:-}" ]]; then
        log::error "Cannot run flock::run while flock::start holds the lock_fd=$lock_fd"
        return 1
    fi

    local rc=0
    flock -x -E 199 -w "$APT_LOCK_TIMEOUT" "$APT_LOCKFILE" "$@" || rc=$?

    case "$rc" in
    0)
        return 0
        ;;
    199)
        log::error "failed to acquire svmkit lock within $APT_LOCK_TIMEOUT seconds while running: $(array::join " " "$@")"
        return 199
        ;;
    *)
        log::error "failed to get svmkit lock with exit code $rc while running: $(array::join " " "$@")"
        return "$rc"
        ;;
    esac
}

svmkit::apt::get() {
    log::info "Acquiring svmkit lock and running apt-get..."

    DEBIAN_FRONTEND=noninteractive svmkit::flock::run sudo -E apt-get -qy \
        -o APT::Lock::Timeout="$APT_LOCK_TIMEOUT" \
        -o DPkg::Lock::Timeout="$APT_LOCK_TIMEOUT" \
        "$@"
}

svmkit::apt::update() {
    svmkit::apt::get update
}

apt::setup-abk-apt-source() {
    svmkit::apt::update
    svmkit::apt::get install curl gnupg

    svmkit::flock::start
    local did_add_key=false
    if ! grep -q "^deb .*/svmkit dev main" /etc/apt/sources.list /etc/apt/sources.list.d/*; then
        curl -s https://apt.abklabs.com/keys/abklabs-archive-dev.asc | svmkit::sudo apt-key add -
        echo "deb https://apt.abklabs.com/svmkit dev main" |
            svmkit::sudo tee /etc/apt/sources.list.d/svmkit.list >/dev/null
        did_add_key=true
    fi
    svmkit::flock::end

    if $did_add_key; then
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

    svmkit::flock::start
    id sol >/dev/null 2>&1 || svmkit::sudo adduser --disabled-password --gecos "" sol
    svmkit::sudo mkdir -p "/home/sol"
    svmkit::sudo chown -f -R sol:sol "/home/sol"
    svmkit::sudo chmod 750 "/home/sol"

    username=$(whoami)
    id -nGz "$username" | grep -qzxF sol || svmkit::sudo adduser "$username" sol

    cat <<EOF | svmkit::sudo tee /etc/security/limits.d/50-sol.conf >/dev/null
sol    soft    nofile    1000000
sol    hard    nofile    1000000
EOF

    svmkit::sudo chown root:root /etc/security/limits.d/50-sol.conf
    svmkit::sudo chmod 644 /etc/security/limits.d/50-sol.conf

    svmkit::flock::end
}
