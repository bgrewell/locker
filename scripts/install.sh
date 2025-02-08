#!/bin/bash
set -euo pipefail

# Define installation directories.
INSTALL_DIR="/opt/locker"
BIN_DIR="${INSTALL_DIR}/bin"
PAM_DIR="${INSTALL_DIR}/pam"
CONFIG_DIR="/etc/locker"
SYSTEMD_DIR="/etc/systemd/system"

echo "Creating installation directories..."
sudo mkdir -p "${BIN_DIR}"
sudo mkdir -p "${PAM_DIR}"
sudo mkdir -p "${CONFIG_DIR}"

# Download the latest release tarball from GitHub.
# This assumes the asset is available at this URL.
echo "Downloading latest release from GitHub..."
curl -L -o locker.tar.gz https://github.com/bgrewell/locker/releases/latest/download/locker.tar.gz

# Create a temporary directory for extraction.
TMPDIR=$(mktemp -d)
echo "Extracting files to temporary directory..."
tar -xzf locker.tar.gz -C "$TMPDIR"

# Install binaries to /opt/locker/bin.
echo "Installing binaries to ${BIN_DIR}..."
sudo cp "$TMPDIR/locker" "${BIN_DIR}/"
sudo cp "$TMPDIR/lockerd" "${BIN_DIR}/"
sudo chmod +x "${BIN_DIR}/locker" "${BIN_DIR}/lockerd"

# Install PAM module if available.
if [ -f "$TMPDIR/pam_locker.so" ]; then
    echo "Installing PAM module to ${PAM_DIR}..."
    sudo cp "$TMPDIR/pam_locker.so" "${PAM_DIR}/"
    sudo chmod 644 "${PAM_DIR}/pam_locker.so"
fi

# Install configuration file if not already present.
if [ -f "$TMPDIR/config/config.yaml" ]; then
    if [ ! -f "${CONFIG_DIR}/config.yaml" ]; then
        echo "Installing default configuration to ${CONFIG_DIR}..."
        sudo cp "$TMPDIR/config/config.yaml" "${CONFIG_DIR}/"
    else
        echo "Configuration file already exists in ${CONFIG_DIR}, skipping."
    fi
else
    echo "[WARN] config.yaml not found in the release archive."
fi

# Install the systemd service file.
if [ -f "$TMPDIR/service/lockerd.service" ]; then
    echo "Installing systemd service file..."
    sudo cp "$TMPDIR/service/lockerd.service" "${SYSTEMD_DIR}/"
    sudo chmod 644 "${SYSTEMD_DIR}/lockerd.service"
else
    echo "[WARN] lockerd.service not found in the release archive."
fi

echo "Cleaning up temporary files..."
rm -rf "$TMPDIR" locker.tar.gz

echo "Reloading systemd daemon..."
sudo systemctl daemon-reload

# --- PAM Module Insertion ---
# Update /etc/pam.d/sshd to insert the pam_locker module.

PAM_FILE="/etc/pam.d/sshd"
BACKUP_FILE="/etc/pam.d/sshd.bak"
MODULE_ENTRY="pam_locker.so"

if [ -f "${PAM_DIR}/pam_locker.so" ]; then
    echo "Updating PAM configuration in ${PAM_FILE}..."

    # Create a backup if it doesn't exist.
    if [ ! -f "$BACKUP_FILE" ]; then
        echo "Backing up ${PAM_FILE} to ${BACKUP_FILE}"
        sudo cp "$PAM_FILE" "$BACKUP_FILE"
    fi

    # Insert the auth module entry before the @include common-auth line.
    if ! grep -q "^auth[[:space:]]\+required[[:space:]]\+$MODULE_ENTRY" "$PAM_FILE"; then
        echo "Inserting auth module entry..."
        sudo sed -i '/^@include[[:space:]]\+common-auth/i auth    required    '"$MODULE_ENTRY" "$PAM_FILE"
    fi

    # Insert the account module entry immediately after the line containing pam_nologin.so.
    if ! grep -q "^account[[:space:]]\+required[[:space:]]\+$MODULE_ENTRY" "$PAM_FILE"; then
        echo "Inserting account module entry..."
        sudo sed -i '/pam_nologin.so/a account    required    '"$MODULE_ENTRY" "$PAM_FILE"
    fi

    # Insert the session module entry before the line that includes pam_selinux.so open.
    if ! grep -q "^session[[:space:]]\+required[[:space:]]\+$MODULE_ENTRY" "$PAM_FILE"; then
        echo "Inserting session module entry..."
        sudo sed -i '/pam_selinux.so[[:space:]]\+open/i session    required    '"$MODULE_ENTRY" "$PAM_FILE"
    fi

    echo "PAM configuration updated in ${PAM_FILE}."
fi
# --- End PAM Module Insertion ---

# Automatically enable and start the systemd service.
echo "Enabling and starting lockerd.service..."
sudo systemctl enable lockerd.service
sudo systemctl start lockerd.service

# Create a symlink for the PAM module.
# Determine the system's PAM module directory.
if [ -d "/lib/security" ]; then
    PAM_TARGET="/lib/security"
elif [ -d "/lib/x86_64-linux-gnu/security/" ]; then
    PAM_TARGET="/lib/x86_64-linux-gnu/security/"
else
    echo "[WARN] PAM module directory not found; skipping linking PAM module."
    PAM_TARGET=""
fi

if [ -n "$PAM_TARGET" ] && [ -f "${PAM_DIR}/pam_locker.so" ]; then
    echo "Creating symlink for PAM module in ${PAM_TARGET}..."
    sudo ln -sf "${PAM_DIR}/pam_locker.so" "${PAM_TARGET}/pam_locker.so"
fi

# Create wrapper scripts for "lock" and "unlock" commands.
echo "Creating wrapper script for 'lock'..."
sudo tee /usr/local/bin/lock > /dev/null <<'EOF'
#!/bin/bash
/opt/locker/bin/locker lock "$@"
EOF
sudo chmod +x /usr/local/bin/lock

echo "Creating wrapper script for 'unlock'..."
sudo tee /usr/local/bin/unlock > /dev/null <<'EOF'
#!/bin/bash
/opt/locker/bin/locker unlock "$@"
EOF
sudo chmod +x /usr/local/bin/unlock

# Create a symlink for "locker" itself.
echo "Creating symlink for 'locker'..."
sudo ln -sf /opt/locker/bin/locker /usr/local/bin/locker

echo "Installation complete."
echo "You can now use the following commands:"
echo "  sudo systemctl status lockerd.service   - To check the daemon status."
echo "  /usr/local/bin/lock                      - To lock the system."
echo "  /usr/local/bin/unlock                    - To unlock the system."
echo "  /usr/local/bin/locker                    - To run the locker CLI."
