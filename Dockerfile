# Dockerfile for cross-platform Pomo builds
FROM golang:1.12-stretch

ENV OSXCROSS_REPO="https://github.com/tpoechtrager/osxcross.git"
ENV OSX_SDK_TARBALL="https://s3.dockerproject.org/darwin/v2/MacOSX10.11.sdk.tar.xz"

RUN \
    apt-get update -yyq \
    && apt-get install -yyq \
        clang \
        libxml2 \
        patch \
        xz-utils

RUN \
    mkdir /build \
    && cd /build \
    && git clone --depth=1 "$OSXCROSS_REPO" \
    && cd osxcross/tarballs \
    && wget "$OSX_SDK_TARBALL" \
    && cd .. \
    && UNATTENDED=1 ./build.sh

ENV PATH="$PATH:/build/osxcross/target/bin"
