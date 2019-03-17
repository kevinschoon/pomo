#!/bin/bash

POMO_VERSION="0.7.1"
OSX_MD5="369a489fc1e9af234cd5db099c3efd83"
LINUX_MD5="973f5c83218d1d3df5e43a2e6c171b7f"
OSX_TARBALL="https://github.com/kevinschoon/pomo/releases/download/$POMO_VERSION/pomo-$POMO_VERSION-darwin-amd64"
LINUX_TARBALL="https://github.com/kevinschoon/pomo/releases/download/$POMO_VERSION/pomo-$POMO_VERSION-linux-amd64"

install_pomo() {
    echo "Installing Pomo..."
    if [[ "$OSTYPE" == darwin* ]] ; then {
        curl -L -o pomo "$OSX_TARBALL" && \
        [[ $(md5 -r pomo) == "$OSX_MD5" ]] && \
        chmod +x pomo && \
        ./pomo -v
    } elif [[ "$OSTYPE" == linux* ]] ; then {
        curl -L -o pomo "$LINUX_TARBALL" && \
        echo "$LINUX_MD5 pomo" | md5sum -c - && \
        chmod +x pomo && \
        ./pomo -v
    } else {
        echo "cannot detect OS type"
        return 1
    }
    fi
    echo "Pomo $POMO_VERSION installed, copy ./pomo to somewhere on your path."
}

install_pomo
