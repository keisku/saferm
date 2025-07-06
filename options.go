package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Options represents the options for saferm.
type Options struct {
	Args     []string
	TrashDir string

	// Command line options
	Force     bool
	Verbose   bool
	Recursive bool
	Dir       bool
}

// Run runs saferm.
func (o *Options) Run() error {
	for _, target := range o.Args {
		base := filepath.Base(target)
		if base == "." || base == ".." {
			return fmt.Errorf("refusing to remove '%s'", target)
		}
		absTarget, err := filepath.Abs(target)
		if err == nil && absTarget == "/" {
			return fmt.Errorf("refusing to operate on '/'")
		}

		info, err := os.Lstat(target)
		if err != nil {
			if !o.Force {
				return fmt.Errorf("failed to move '%s' to Trash bin: %v", target, err)
			}
			return nil
		}

		if info.IsDir() {
			if o.Recursive {
				// allow recursive directory removal
			} else if o.Dir {
				// allow only if directory is empty
				empty, err := o.isDirEmpty(target)
				if err != nil {
					return fmt.Errorf("cannot check if directory '%s' is empty: %v", target, err)
				}
				if !empty {
					return fmt.Errorf("cannot move directory '%s' to Trash bin: directory not empty", target)
				}
			} else {
				return fmt.Errorf("cannot move directory '%s' to Trash bin: is a directory (use -r or -R to move recursively, or -d to remove empty directories)", target)
			}
		}

		dest := filepath.Join(o.TrashDir, base)
		err = os.Rename(target, dest)
		if err != nil && !o.Force {
			return fmt.Errorf("failed to move '%s' to Trash bin: %v", target, err)
		}
	}
	return nil
}

func (o *Options) isDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// Read at most one entry
	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
