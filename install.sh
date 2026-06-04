#!/bin/bash
set -e

echo "============================================="
echo "        Devopin Global Installer             "
echo "============================================="

# Detect OS
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
if [ "$OS" != "linux" ] && [ "$OS" != "darwin" ]; then
    echo "Error: OS $OS is not supported by this installer."
    exit 1
fi

# Detect Architecture
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64|amd64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Error: Architecture $ARCH is not supported."
        exit 1
        ;;
esac

echo "Detected System: $OS ($ARCH)"

# Fetch latest release version from GitHub
echo "Fetching latest version..."
LATEST_TAG=$(curl -s "https://api.github.com/repos/gabutlabs/godevopin/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_TAG" ]; then
    # Fallback to the known beta tag if API limit is reached or repository is private
    LATEST_TAG="v0.2.0-beta"
    echo "Could not fetch latest release automatically, defaulting to $LATEST_TAG"
else
    echo "Latest version found: $LATEST_TAG"
fi

BINARY_NAME="godevopin-${OS}-${ARCH}"
DOWNLOAD_URL="https://github.com/gabutlabs/godevopin/releases/download/${LATEST_TAG}/${BINARY_NAME}"

echo "Downloading $BINARY_NAME from GitHub Releases..."
TMP_DIR=$(mktemp -d)
if ! curl -f -L -o "$TMP_DIR/devopin" "$DOWNLOAD_URL"; then
    echo "Error: Binary not found on GitHub Releases for this OS/Arch, or download failed."
    echo "URL: $DOWNLOAD_URL"
    echo "If your repository is Private, you cannot use this installer directly without authentication."
    rm -rf "$TMP_DIR"
    exit 1
fi

# A GitHub "Not Found" page might occasionally be returned with HTTP 200, though curl -f usually catches 404s.
# To be safe, we verify the file size. A valid Go binary is > 1MB, while an error page is tiny.
FILE_SIZE=$(wc -c < "$TMP_DIR/devopin" | tr -d ' ')
if [ "$FILE_SIZE" -lt 10000 ]; then
    echo "Error: Downloaded file is too small to be a valid binary. It might be a 'Not Found' page."
    echo "URL: $DOWNLOAD_URL"
    echo "If your repository is Private, you cannot use this installer directly without authentication."
    rm -rf "$TMP_DIR"
    exit 1
fi

chmod +x "$TMP_DIR/devopin"

echo "Installing binary to /usr/local/bin/devopin (may require sudo password)..."
sudo mv "$TMP_DIR/devopin" /usr/local/bin/devopin

echo "Setting up configuration directory at /opt/devopin..."
sudo mkdir -p /opt/devopin

CONFIG_FILE="/opt/devopin/config.yaml"
if [ ! -f "$CONFIG_FILE" ]; then
    echo "Downloading example configuration..."
    sudo curl -sSL -o "$CONFIG_FILE" "https://raw.githubusercontent.com/gabutlabs/godevopin/main/configs/config.yaml.example"
    
    # If the file couldn't be downloaded from main (e.g. repo is private), create a basic one
    if grep -q "404: Not Found" "$CONFIG_FILE" || [ ! -s "$CONFIG_FILE" ]; then
        echo "Could not download config.yaml.example, generating a default one..."
        sudo bash -c 'cat > /opt/devopin/config.yaml <<EOF
database:
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: "password"
  dbname: "devopin"
app:
  jwt_secret: "change-this-secret-in-production"
EOF'
    fi
    echo "Configuration created at $CONFIG_FILE"
else
    echo "Configuration already exists at $CONFIG_FILE, skipping..."
fi

rm -rf "$TMP_DIR"

# Setup systemd services if on Linux
if [ "$OS" = "linux" ] && command -v systemctl >/dev/null 2>&1; then
    echo "Setting up systemd services for devopin-serve and devopin-worker..."
    
    sudo bash -c 'cat > /etc/systemd/system/devopin-serve.service <<EOF
[Unit]
Description=Devopin Serve Service
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/devopin serve
WorkingDirectory=/opt/devopin
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF'

    sudo bash -c 'cat > /etc/systemd/system/devopin-worker.service <<EOF
[Unit]
Description=Devopin Worker Service
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/devopin worker
WorkingDirectory=/opt/devopin
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF'

    sudo systemctl daemon-reload
    sudo systemctl enable devopin-serve devopin-worker
    sudo systemctl start devopin-serve devopin-worker || true
    echo "Systemd services (devopin-serve, devopin-worker) created, enabled, and started."
fi

echo "============================================="
echo "Installation Successful!"
echo "You can now run 'devopin serve' or 'devopin worker' from anywhere."
if [ "$OS" = "linux" ] && command -v systemctl >/dev/null 2>&1; then
    echo "Systemd services are configured and running. Manage them with:"
    echo "  sudo systemctl status devopin-serve"
    echo "  sudo systemctl status devopin-worker"
fi
echo "Don't forget to edit /opt/devopin/config.yaml with your settings."
echo "============================================="
