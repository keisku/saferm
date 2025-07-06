# saferm

A safer alternative to `rm`.
Instead of permanently deleting files and directories, `saferm` moves them to your `~/.Trash` directory, making accidental deletions easy to recover.

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
