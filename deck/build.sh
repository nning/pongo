#!/bin/sh
# See https://ebiten.org/blog/steam.html#Linux

name=pongo
STEAM_RUNTIME_VERSION=0.20210817.0
GO_VERSION=$(go env GOVERSION)

mkdir -p .cache/${STEAM_RUNTIME_VERSION}

# Download binaries for amd64.
if [[ ! -f .cache/${STEAM_RUNTIME_VERSION}/com.valvesoftware.SteamRuntime.Sdk-amd64,i386-scout-sysroot.Dockerfile ]]; then
    (cd .cache/${STEAM_RUNTIME_VERSION}; curl --location --remote-name https://repo.steampowered.com/steamrt-images-scout/snapshots/${STEAM_RUNTIME_VERSION}/com.valvesoftware.SteamRuntime.Sdk-amd64,i386-scout-sysroot.Dockerfile)
fi
if [[ ! -f .cache/${STEAM_RUNTIME_VERSION}/com.valvesoftware.SteamRuntime.Sdk-amd64,i386-scout-sysroot.tar.gz ]]; then
    (cd .cache/${STEAM_RUNTIME_VERSION}; curl --location --remote-name https://repo.steampowered.com/steamrt-images-scout/snapshots/${STEAM_RUNTIME_VERSION}/com.valvesoftware.SteamRuntime.Sdk-amd64,i386-scout-sysroot.tar.gz)
fi
if [[ ! -f .cache/${GO_VERSION}.linux-amd64.tar.gz ]]; then
    (cd .cache; curl --location --remote-name https://golang.org/dl/${GO_VERSION}.linux-amd64.tar.gz)
fi

# Build for amd64.
(cd .cache/${STEAM_RUNTIME_VERSION}; docker build -f com.valvesoftware.SteamRuntime.Sdk-amd64,i386-scout-sysroot.Dockerfile -t steamrt_scout_amd64:latest .)
docker run --rm --workdir=/work --volume $(pwd):/work:z steamrt_scout_amd64:latest /bin/sh -c "
export PATH=\$PATH:/usr/local/go/bin
export CGO_CFLAGS=-std=gnu99

rm -rf /usr/local/go && tar -C /usr/local -xzf .cache/${GO_VERSION}.linux-amd64.tar.gz

make clean build
"
# go build -o ${name}_linux_amd64 .
