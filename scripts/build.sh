#!/bin/bash
set -e

# Get project root
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "---Building Web Assets ---"
cd "$ROOT/web"
bun run build

echo "--- Syncing Assets to Go ---"
rm -rf "$ROOT/ui/public"
mkdir -p "$ROOT/ui/public"
cp -r "$ROOT/web/out/"* "$ROOT/ui/public/"

echo "--- Done! Created ./sector-one ---"