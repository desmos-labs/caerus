#!/usr/bin/env bash

SCRIPT_NAME=$(basename "$0")
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
DESMOS_BIN="$SCRIPT_DIR/desmos"

# Remove the previous version of desmos if exists
if test -f "$SCRIPT_DIR/desmos"; then
    # rm "$DESMOS_BIN"
    # TODO: Allow to provide a --force flag that allows overwriting the binary.
    # Instead of removing the binary, use the one we have
    exit 0
fi

if [ -z "$1" ]; then
  echo "Please provide the Desmos version. E.g:";
  echo "$SCRIPT_NAME 4.4.1";
fi

# Set version variable from first argument
VERSION=$1

# Build desmos download url from the received arg
DESMOS_URL="https://github.com/desmos-labs/desmos/releases/download/v$VERSION/desmos-$VERSION-linux-amd64"

# Force the script to exit at the first error
set -e

echo "Downloading desmos version: $VERSION"

# Download the binary
wget --show-progress -O "$DESMOS_BIN" "$DESMOS_URL"

# Make the binary executable
chmod +x "$DESMOS_BIN"
