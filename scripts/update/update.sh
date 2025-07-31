#!/usr/bin/env bash
set -euo pipefail

# resolve project root (one level up from this script)
ROOT_DIR="$(cd "$(dirname ../../)/.." && pwd)"

# iterate over each subdirectory in project root
for dir in "$ROOT_DIR"/*/; do
  # skip non-dirs just in case
  [ -d "$dir" ] || continue

  # get bare folder name, without trailing slash
  name=${dir##*/}        # e.g. "/path/to/vendor/" → "vendor"
  name=${name%/}

  # skip specific folders
  if [[ "$name" == "scripts" || "$name" == "third_party" ]]; then
    echo "→ Skipping $name"
    continue
  fi

  echo "→ Updating Go modules in $name"
  (
    cd "$dir"
    go get -u ./...
  )
done