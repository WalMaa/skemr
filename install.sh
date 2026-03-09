#!/bin/sh

set -eu

# If running as root, install to /usr/local/bin, otherwise install to ~/.local/bin
if [ "$(id -u)" -eq 0 ]; then
  INSTALL_DIR="/usr/local/bin"
else
  INSTALL_DIR="${HOME}/.local/bin"
  # Create the directory if it doesn't exist
  mkdir -p "$INSTALL_DIR"
fi

URL="https://github.com/WalMaa/skemr/releases/download/v0.1.2/skemr-cli_Linux_x86_64"

curl --request GET -fL \
     --url "$URL"\
     --output "$INSTALL_DIR/skemr"

chmod +x "$INSTALL_DIR/skemr"

case ":$PATH:" in
  *":$INSTALL_DIR:"*)
    echo "Installed to $INSTALL_DIR"
    echo "Run: skemr"
    ;;
  *)
    echo "Installed to $INSTALL_DIR"
    echo "$INSTALL_DIR is not on your PATH."
    echo "For the current shell session, run:"
    echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
    echo "To make it permanent, add this to your shell config:"
    echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
    ;;
esac