#!/usr/bin/env bash
set -euo pipefail

ALL_RIDS=(
    linux-x64 linux-arm64 linux-musl-x64 linux-musl-arm64
    win-x64 win-arm64
    osx-x64 osx-arm64
)

SILENT=false
RIDS=()
VERSION=""
CONFIG="Release"

usage() {
    cat <<EOF
Использование: ./build.sh [платформа...] [опции]

Платформы:
$(printf "  %s\n" "${ALL_RIDS[@]}")

Опции:
  --silent, -s  Без вывода логов
  --version, -v Указать версию (по умолч. git tag или дата)
  --help, -h    Эта справка

Если платформы не указаны — собираются все.

Примеры:
  ./build.sh linux-x64 win-x64
  ./build.sh --silent
  ./build.sh osx-arm64 --version 1.2.3
  ./build.sh --help
EOF
    exit 0
}

while [[ $# -gt 0 ]]; do
    case "$1" in
        --help|-h) usage ;;
        --silent|-s) SILENT=true; shift ;;
        --version|-v) VERSION="$2"; shift 2 ;;
        --*) echo "Неизвестный флаг: $1"; exit 1 ;;
        *) RIDS+=("$1"); shift ;;
    esac
done

[[ ${#RIDS[@]} -eq 0 ]] && RIDS=("${ALL_RIDS[@]}")
[[ -z "$VERSION" ]] && VERSION="$(git describe --tags --always 2>/dev/null || date +%Y%m%d)"

OUT="dist/$VERSION"
VERBOSITY=""
$SILENT && VERBOSITY="--verbosity quiet"

echo "=== Cubenet Launcher v$VERSION ==="
echo "Папка: $OUT"
echo

for rid in "${RIDS[@]}"; do
    $SILENT || echo "→ $rid"
    dotnet publish src/CubenetLauncher \
        -c "$CONFIG" \
        -r "$rid" \
        --self-contained true \
        -o "$OUT/$rid" $VERBOSITY
done

echo
echo "=== Готово ==="
ls -lh "$OUT"/*/Cubenet_Launcher* 2>/dev/null
