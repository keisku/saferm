# saferm

![GitHub release (latest by date)](https://img.shields.io/github/v/release/keisku/saferm?style=flat-square)

![Trash icon](./trash.png)

A safer alternative to `rm`.
Instead of permanently deleting files and directories, `saferm` moves them to your `~/.Trash` directory, making accidental deletions easy to recover.

## Installation

Go to [Releases](https://github.com/keisku/saferm/releases) and download the binary for your platform.

Example for macOS:

```
curl -L "https://github.com/keisku/saferm/releases/latest/download/saferm-darwin-amd64" -o /usr/local/bin/saferm && sudo chmod +x /usr/local/bin/saferm
```

## Usage

Remove a file (moves to Trash):

```
saferm myfile.txt
```

Remove an empty directory:

```
saferm -d mydir
```

Remove a directory and all its contents:

```
saferm -r mydir
```

Remove multiple files:

```
saferm file1.txt file2.txt
```

Force remove (ignore errors for missing files):

```
saferm -f missing.txt
```

## Safety Notes

- Files and directories are moved to `~/.Trash` and can be recovered manually.
- `saferm` will not remove `.` or `..` or operate on `/` for safety.
- Not all systems have `~/.Trash` by default; `saferm` will create it if needed.
