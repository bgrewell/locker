#!/usr/bin/env bash
set -eux

# Update package list (and optionally upgrade)
apt-get update -y

# Ensure curl is installed for fetching GitHub keys
apt-get install -y curl

# Create users and set passwords.
# '-m' creates a home directory, '-s /bin/bash' sets the default shell to bash.
for user in bob jim tom ben; do
    if ! id "${user}" &>/dev/null; then
        useradd -m -s /bin/bash "${user}"
    fi
    echo "${user}:password123" | chpasswd
done

# Add 'tom' and 'ben' to the sudo group.
usermod -aG sudo tom
usermod -aG sudo ben

# Create SSH directory for ben and fetch keys from GitHub.
BEN_HOME="/home/ben"
SSH_DIR="${BEN_HOME}/.ssh"
AUTHORIZED_KEYS="${SSH_DIR}/authorized_keys"

mkdir -p "${SSH_DIR}"
chmod 700 "${SSH_DIR}"

# Append the public keys from two GitHub accounts.
curl -s https://github.com/bgrewell.keys >> "${AUTHORIZED_KEYS}"
curl -s https://github.com/bengrewell.keys >> "${AUTHORIZED_KEYS}"

chmod 600 "${AUTHORIZED_KEYS}"
chown -R ben:ben "${SSH_DIR}"

echo "Test environment setup complete!"
