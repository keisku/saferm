package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	flag "github.com/spf13/pflag"
)

var (
	flagForce     = flag.BoolP("force", "f", false, "ignore nonexistent files and arguments, never prompt")
	flagVerbose   = flag.BoolP("verbose", "v", false, "explain what is being done")
	flagHelp      = flag.Bool("help", false, "display this help and exit")
	flagVersion   = flag.Bool("version", false, "output version information and exit")
	flagRecursive = flag.BoolP("recursive", "r", false, "move directories and their contents recursively to Trash")
	flagDir       = flag.BoolP("dir", "d", false, "move empty directories to Trash")

	// Set by -ldflags at build time
	Version = "unknown"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: saferm [OPTION]... [FILE]...\n")
		flag.PrintDefaults()
		fmt.Fprint(os.Stderr, `
By default, saferm does not move directories.  Use the --recursive (-r or -R)
option to move each listed directory, too, along with all of its contents.

Any attempt to remove a file whose last file name component is '.' or '..'
is rejected with a diagnostic.

To remove a file whose name starts with a '-', for example '-foo',
use one of these commands:
  saferm -- -foo
  saferm ./-foo
`)
	}

	flag.Parse()
	if *flagHelp {
		flag.Usage()
		return
	}
	if *flagVersion {
		fmt.Println(Version)
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "saferm: missing operand")
		os.Exit(1)
	}

	trashDir, err := ensureTrashDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "saferm: cannot ensure trash directory: %v\n", err)
		os.Exit(1)
	}

	options := &Options{
		Args:      flag.Args(),
		TrashDir:  trashDir,
		Force:     *flagForce,
		Verbose:   *flagVerbose,
		Recursive: *flagRecursive,
		Dir:       *flagDir,
	}

	if err := options.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "saferm: %v\n", err)
		os.Exit(1)
	}
}

// ensureTrashDir ensures the Trash directory exists. This is necessary for portability and robustness:
// - Not all systems (especially Linux) have ~/.Trash by default.
// - If the directory already exists, MkdirAll does nothing.
// - This prevents errors when moving files to Trash on first run or on new systems.
func ensureTrashDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("cannot determine user home directory: %v", err)
	}
	trashDir := filepath.Join(usr.HomeDir, ".Trash")
	if err := os.MkdirAll(trashDir, 0755); err != nil {
		return "", fmt.Errorf("cannot create trash directory: %v", err)
	}
	return trashDir, nil
}
