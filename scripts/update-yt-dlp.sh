#!/bin/sh
set -e

NEW_FILE="/tmp/yt-dlp.new"
CURRENT_FILE="/usr/local/bin/yt-dlp"

echo "$(date): Update yt-dlp started ..."
if wget -q -O "$NEW_FILE" https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp; then
    chmod +x "$NEW_FILE"
    if "$NEW_FILE" --version; then
        mv "$NEW_FILE" "$CURRENT_FILE"
        echo "$(date): Update yt-dlp finished successfully"
        exit 0
    fi
fi

echo "$(date): Update yt-dlp failed"
rm -f "$NEW_FILE"
exit 1
