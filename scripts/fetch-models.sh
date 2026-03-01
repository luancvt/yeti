#!/usr/bin/env bash
set -euo pipefail

if [ $# -ne 2 ]; then
    echo "Usage: fetch-models.sh <name> <repo-path>"
    echo "Example: fetch-models.sh xr-7112 vendor/cisco/xr/7112"
    exit 1
fi

name="$1"
repo_path="$2"

if [ -d "models/$name" ]; then
    echo "models/$name already exists. Remove it first to re-download."
    exit 1
fi

tmpdir=$(mktemp -d)
trap "rm -rf $tmpdir" EXIT

echo "Fetching $repo_path from YangModels/yang..."
git clone --quiet --depth 1 --filter=blob:none --sparse \
    https://github.com/YangModels/yang.git "$tmpdir/yang"
git -C "$tmpdir/yang" sparse-checkout set "$repo_path"
cp -r "$tmpdir/yang/$repo_path" "models/$name"
echo "Downloaded to models/$name"
