#!/bin/bash
set -e

# Get project root
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "--- Building Web Assets ---"
cd "$ROOT/web"
bun run build

echo "--- Syncing Assets to Go ---"
rm -rf "$ROOT/src/ui/public"
mkdir -p "$ROOT/src/ui/public"
cp -r "$ROOT/web/out/"* "$ROOT/src/ui/public/"

echo "--- Building Go Engine ---"
cd "$ROOT"
go build -o sector-one ./src/cmd/engine

echo "--- Done! Created ./sector-one ---"
