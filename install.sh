#!/bin/bash

POMO_VERSION="0.7.0"
OSX_MD5="4aa452e021d07bbb2c7c3a15839fe5df"
LINUX_MD5="84551ab258902e62d9d8e14368bfdf4b"
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
