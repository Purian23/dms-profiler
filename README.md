# DMS Profiler

Linux Chrome Profile: when set as the default browser, it opens **Google Chrome Stable** URLs in the right **profile** using `--profile-directory`, driven by a TOML config.

## Install

```bash
make install
```

This single command builds the binary, installs it to `/usr/local/bin/` (requires sudo), copies the desktop entry to `~/.local/share/applications/`, creates `~/.config/dms-profiler/config.toml` from the example on first install (never overwrites an existing config), and registers dms-profiler as the default handler for `http://` and `https://`. Existing desktop entry files are backed up as `.bak` before being replaced.

To install to a different path:

```bash
make install INSTALL_BIN=~/.local/bin/dms-profiler
```

## Configure

Edit `~/.config/dms-profiler/config.toml` after the first install. Profile names must match what Chrome shows. To list your profiles:

```bash
python3 -c "
import json, os
with open(os.path.expanduser('~/.config/google-chrome/Local State')) as f:
    data = json.load(f)
for folder, info in sorted(data['profile']['info_cache'].items()):
    print(f'{folder:20} -> {info[\"name\"]}')
"
```

Each `[[rules]]` entry uses either:

- **`match`** — one URL prefix string, or  
- **`matches`** — an array of prefixes for the **same** `profile` (less repetition).

Rules are evaluated **top to bottom**; the **first** prefix that matches the URL wins. See `config.example.toml`.

Profile values:

- **Folder names** (`Default`, `Profile 1`, …) are used as-is.
- **Names shown in Chrome** (`Work`, `Personal`, …) are resolved from `Local State` under `user_data_dir`.

## Re-registering the handler

Fedora/GNOME updates or Chrome's "set as default browser" prompt can reset the MIME handler. Fix it with:

```bash
make register-handler
```

## Uninstall

```bash
make uninstall
```

Removes the binary and desktop entry and clears the MIME associations. Your config at `~/.config/dms-profiler/` is **not** removed.

## CLI

```text
dms-profiler [flags] <url>

  -config string
        path to config.toml (default ~/.config/dms-profiler/config.toml)
  -dry-run
        print the command line and exit (does not launch Chrome)
  -print-cmd
        same as -dry-run
```

Example:

```bash
dms-profiler -dry-run 'https://example.com/'
```

## Scope

Chrome Stable, single global `user_data_dir`, prefix URL matching (no mid-URL wildcards).
