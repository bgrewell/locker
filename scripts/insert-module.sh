#!/bin/bash
set -euo pipefail

PAM_FILE="/etc/pam.d/sshd"
BACKUP_FILE="/etc/pam.d/sshd.bak"

# Create a backup if it doesn't exist.
if [ ! -f "$BACKUP_FILE" ]; then
    echo "Backing up $PAM_FILE to $BACKUP_FILE"
    cp "$PAM_FILE" "$BACKUP_FILE"
fi

# Define the module entry.
MODULE_ENTRY="pam_locker.so"

# Insert the auth module entry before the @include common-auth line.
if ! grep -q "^auth[[:space:]]\+required[[:space:]]\+$MODULE_ENTRY" "$PAM_FILE"; then
    echo "Inserting auth module entry..."
    sed -i '/^@include[[:space:]]\+common-auth/i auth    required    '"$MODULE_ENTRY" "$PAM_FILE"
fi

# Insert the account module entry immediately after the line containing pam_nologin.so.
if ! grep -q "^account[[:space:]]\+required[[:space:]]\+$MODULE_ENTRY" "$PAM_FILE"; then
    echo "Inserting account module entry..."
    sed -i '/pam_nologin.so/a account    required    '"$MODULE_ENTRY" "$PAM_FILE"
fi

# Insert the session module entry before the line that includes pam_selinux.so open.
if ! grep -q "^session[[:space:]]\+required[[:space:]]\+$MODULE_ENTRY" "$PAM_FILE"; then
    echo "Inserting session module entry..."
    sed -i '/pam_selinux.so[[:space:]]\+open/i session    required    '"$MODULE_ENTRY" "$PAM_FILE"
fi

echo "PAM configuration updated in $PAM_FILE."
