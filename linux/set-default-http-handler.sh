#!/usr/bin/env bash
# Re-register dms-profiler as the handler for http(s) links.
# Run this after Fedora/GNOME updates or if a "default browser" prompt reset Chrome
# as the handler (xdg-mime associations live in ~/.config/mimeapps.list).

set -euo pipefail

DESKTOP_ID="io.github.dms-profiler.desktop"
USER_APPS="${XDG_DATA_HOME:-$HOME/.local/share}/applications"

if [[ ! -f "$USER_APPS/$DESKTOP_ID" ]]; then
	echo "Warning: $USER_APPS/$DESKTOP_ID not found." >&2
	echo "Install it first, e.g.: cp linux/$DESKTOP_ID \"$USER_APPS/\"" >&2
fi

xdg-mime default "$DESKTOP_ID" x-scheme-handler/http
xdg-mime default "$DESKTOP_ID" x-scheme-handler/https

echo "HTTP/HTTPS default handler: $DESKTOP_ID"
echo -n "  http:  "
xdg-mime query default x-scheme-handler/http
echo -n "  https: "
xdg-mime query default x-scheme-handler/https
