#!/bin/bash

# install dependencies and common tools
apk add wget bmon vim htop
apk add gcc glib libatomic libcap2 libgcc libpcap libpcap-dev musl-dev zstd-libs zlib

# install golang
wget -O go.tgz 'https://dl.google.com/go/go1.21.6.linux-amd64.tar.gz'
tar -C /usr/local -xzf go.tgz
rm go.tgz
export PATH=$PATH:/usr/local/go/bin

# build and install WireMeter
/usr/local/go/bin/go build
mv WireMeter /opt/wiremeter.bin
cp wiremeter.init.d /etc/init.d/wiremeter
rc-update add wiremeter default
rc-service wiremeter start