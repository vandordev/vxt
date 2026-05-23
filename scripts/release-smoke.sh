#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
tmpdir="$(mktemp -d)"

cleanup() {
	rm -rf "$tmpdir"
}
trap cleanup EXIT

cp "$repo_root/smoketest/publicapi/main.go" "$tmpdir/main.go"

cd "$tmpdir"

env GOWORK=off go mod init example.com/vxtsmoke >/dev/null
env GOWORK=off go mod edit -replace "github.com/vandordev/vxt=$repo_root"
env GOWORK=off go get github.com/vandordev/vxt >/dev/null
env GOWORK=off go run .
