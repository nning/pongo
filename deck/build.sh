#!/bin/sh
# See https://ebiten.org/blog/steam.html#Linux

name=pongo
STEAM_RUNTIME_NAME=sniper
STEAM_RUNTIME_VERSION=0.20220718.0

GO_VERSION=$(go env GOVERSION)
FILENAME=com.valvesoftware.SteamRuntime.Sdk-amd64,i386-${STEAM_RUNTIME_NAME}-sysroot

mkdir -p .cache/${STEAM_RUNTIME_VERSION}
cd .cache/${STEAM_RUNTIME_VERSION}
[ ! -f ${FILENAME}.Dockerfile ] &&
	curl --location --remote-name https://repo.steampowered.com/steamrt-images-${STEAM_RUNTIME_NAME}/snapshots/${STEAM_RUNTIME_VERSION}/${FILENAME}.Dockerfile
[ ! -f ${FILENAME}.tar.xz ] &&
	curl --location --remote-name https://repo.steampowered.com/steamrt-images-${STEAM_RUNTIME_NAME}/snapshots/${STEAM_RUNTIME_VERSION}/${FILENAME}.tar.gz

cd ..
[ ! -f ${GO_VERSION}.linux-amd64.tar.gz ] &&
    curl --location --remote-name https://golang.org/dl/${GO_VERSION}.linux-amd64.tar.gz

# Build for amd64.
(cd .cache/${STEAM_RUNTIME_VERSION}; docker build -f ${FILENAME}.Dockerfile -t steamrt_${STEAM_RUNTIME_NAME}_amd64:latest .)
docker run --rm --workdir=/work --volume $(pwd):/work:z steamrt_${STEAM_RUNTIME_NAME}_amd64:latest /bin/sh -c "
export PATH=\$PATH:/usr/local/go/bin
export CGO_CFLAGS=-std=gnu99

rm -rf /usr/local/go && tar -C /usr/local -xzf .cache/${GO_VERSION}.linux-amd64.tar.gz

make clean build
"
