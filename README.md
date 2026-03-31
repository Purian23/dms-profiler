# DMS Profiler

Linux Chrome Profile: when set as the default browser, it opens **Google Chrome Stable** URLs in the right **profile** using `--profile-directory`, driven by a TOML config.

## Build

```bash
go build -o dms-profiler ./cmd/dms-profiler
```

Install the binary somewhere on `PATH`, for example:

```bash
sudo install -m 0755 dms-profiler /usr/local/bin/
```

## Configure

```bash
mkdir -p ~/.config/dms-profiler
cp config.example.toml ~/.config/dms-profiler/config.toml
```

Each `[[rules]]` entry uses either:

- **`match`** — one URL prefix string, or  
- **`matches`** — an array of prefixes for the **same** `profile` (less repetition).

Rules are evaluated **top to bottom**; the **first** prefix that matches the URL wins. See `config.example.toml`.

Profile values:

- **Folder names** (`Default`, `Profile 1`, …) are used as-is.
- **Names shown in Chrome** (`Work`, `Personal`, …) are resolved from `Local State` under `user_data_dir`.

## Desktop integration

Install the desktop entry and point it at your binary path (edit `Exec=` if needed):

```bash
cp linux/io.github.dms-profiler.desktop ~/.local/share/applications/
# Set Exec=/full/path/to/dms-profiler %u if the binary is not on PATH
```

Register as the default handler for HTTP and HTTPS:

```bash
xdg-mime default io.github.dms-profiler.desktop x-scheme-handler/http
xdg-mime default io.github.dms-profiler.desktop x-scheme-handler/https
```

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
