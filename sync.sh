#!/bin/sh
curl -o x86_64-repodata https://repo.voidlinux.eu/current/x86_64-repodata
curl -o i686-repodata https://repo.voidlinux.eu/current/i686-repodata
curl -o armv7l-repodata https://repo.voidlinux.eu/current/armv7l-repodata
curl -o armv6l-repodata https://repo.voidlinux.eu/current/armv6l-repodata


curl -o x86_64-musl-repodata https://repo.voidlinux.eu/current/musl/x86_64-musl-repodata
curl -o armv7l-musl-repodata https://repo.voidlinux.eu/current/musl/armv7l-musl-repodata
curl -o armv6l-musl-repodata https://repo.voidlinux.eu/current/musl/armv6l-musl-repodata

curl -o aarch64-repodata https://repo.voidlinux.eu/current/aarch64/aarch64-repodata
curl -o aarch64-musl-repodata https://repo.voidlinux.eu/current/aarch64/aarch64-musl-repodata
