#!/bin/bash
set -e

echo "============================================="
echo "        Devopin Global Uninstaller           "
echo "============================================="

# Detect OS
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"

if [ "$OS" = "linux" ] && command -v systemctl >/dev/null 2>&1; then
    echo "Stopping and disabling systemd services..."
    sudo systemctl stop devopin-serve devopin-worker || true
    sudo systemctl disable devopin-serve devopin-worker || true
    
    echo "Removing systemd service files..."
    sudo rm -f /etc/systemd/system/devopin-serve.service
    sudo rm -f /etc/systemd/system/devopin-worker.service
    sudo systemctl daemon-reload
fi

echo "Removing binary..."
sudo rm -f /usr/local/bin/devopin

echo "============================================="
echo "Devopin has been successfully uninstalled!"
echo "Note: Configuration files in /opt/devopin were NOT removed."
echo "If you want to remove them completely, run:"
echo "  sudo rm -rf /opt/devopin"
echo "============================================="
