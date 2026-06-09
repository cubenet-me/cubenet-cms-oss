#!/bin/bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"
cd "$ROOT"

usage() {
  cat <<EOF
Usage: $0 [command]

Commands:
  dev       Run in development mode (auto-rebuild on changes)
  build     Build CSS + Templ + Go binary
  docker    Build and start via Docker Compose
  help      Show this help
EOF
}

dev() {
  echo ":: Starting dev mode..."

  # Start tailwind watcher
  (cd backend && npx tailwindcss -i web/static/input.css -o web/static/output.css --watch) &
  TW_PID=$!

  # Start templ watcher
  (cd backend && PATH="$HOME/go/bin:$PATH" templ generate --watch) &
  TM_PID=$!

  cleanup() {
    kill "$TW_PID" "$TM_PID" 2>/dev/null || true
    wait "$TW_PID" "$TM_PID" 2>/dev/null || true
  }
  trap cleanup EXIT INT TERM

  # Run server with reload on recompile
  cd backend && PATH="$HOME/go/bin:$PATH" go run ./cmd/server
}

build() {
  echo ":: Building..."
  cd backend
  npx tailwindcss -i web/static/input.css -o web/static/output.css
  PATH="$HOME/go/bin:$PATH" templ generate
  CGO_ENABLED=0 go build -o server ./cmd/server
  echo ":: Built backend/server"
}

docker_build() {
  echo ":: Building and starting via Docker Compose..."
  docker compose up -d --build --remove-orphans
  echo ":: Done"
}

case "${1:-dev}" in
  dev) dev ;;
  build) build ;;
  docker) docker_build ;;
  help | --help | -h) usage ;;
  *) echo "Unknown command: $1"; usage; exit 1 ;;
esac
