#!/usr/bin/env bash

ENV_FILE=${1}
TARGET_DIR=${2:-$(pwd)}

if [ ! -f "$ENV_FILE" ]; then
    echo "Environment file $ENV_FILE not found"
    exit 1
fi

DIRS=("card" "deck" "repeatro" "sso" "stats")

# Check all directories exist
for dir in "${DIRS[@]}"; do
    if [ ! -d "$TARGET_DIR/$dir" ]; then
        echo "Directory $dir not found in $TARGET_DIR"
        exit 1
    fi
done

# Copy to all directories if all exist
for dir in "${DIRS[@]}"; do
    cp "$ENV_FILE" "$TARGET_DIR/$dir/.env"
done