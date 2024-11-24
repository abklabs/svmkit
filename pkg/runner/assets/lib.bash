apt::setup-abk-apt-source() {
    $APT update
    $APT install curl gnupg
    if ! grep -q "^deb .*/svmkit dev main" /etc/apt/sources.list /etc/apt/sources.list.d/*; then
        curl -s https://apt.abklabs.com/keys/abklabs-archive-dev.asc | $SUDO apt-key add -
        echo "deb https://apt.abklabs.com/svmkit dev main" | $SUDO tee /etc/apt/sources.list.d/svmkit.list >/dev/null
        $APT update
    fi
}

apt::env
