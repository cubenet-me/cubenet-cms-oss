#!/usr/bin/env bash
set -euo pipefail

VERSION=${1:-$(date +%Y%m%d)}
CONFIG=${2:-Release}

echo "=== Публикация Cubenet Launcher v$VERSION ==="

publish() {
    local rid=$1
    local dir="dist/$VERSION/$rid"
    echo "→ $rid"
    dotnet publish src/CubenetLauncher -c "$CONFIG" \
        -r "$rid" \
        --self-contained true \
        -o "$dir"
}

publish "linux-x64"
publish "win-x64"
publish "osx-x64"
publish "osx-arm64"

echo "=== Готово: dist/$VERSION/ ==="
ls -lh "dist/$VERSION"/*/Cubenet_Launcher*
