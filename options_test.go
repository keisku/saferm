package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOptions_Run(t *testing.T) {
	tmpDir := t.TempDir()

	trashDir := filepath.Join(tmpDir, ".Trash")
	os.Mkdir(trashDir, 0755)

	testFile := filepath.Join(tmpDir, "file.txt")
	if err := os.WriteFile(testFile, []byte("data"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	emptyDir := filepath.Join(tmpDir, "emptydir")
	os.Mkdir(emptyDir, 0755)

	nonEmptyDir := filepath.Join(tmpDir, "nonemptydir")
	os.Mkdir(nonEmptyDir, 0755)
	childFile := filepath.Join(nonEmptyDir, "child.txt")
	os.WriteFile(childFile, []byte("child"), 0644)

	tests := []struct {
		name    string
		options Options
		wantErr bool
		check   func() // optional post-check
	}{
		{
			name:    "remove file",
			options: Options{Args: []string{testFile}, TrashDir: trashDir},
			wantErr: false,
			check: func() {
				if _, err := os.Stat(filepath.Join(trashDir, "file.txt")); err != nil {
					t.Errorf("file not moved to trash: %v", err)
				}
			},
		},
		{
			name:    "refuse to remove dot",
			options: Options{Args: []string{"."}, TrashDir: trashDir},
			wantErr: true,
		},
		{
			name:    "refuse to remove dot-dot",
			options: Options{Args: []string{".."}, TrashDir: trashDir},
			wantErr: true,
		},
		{
			name:    "refuse to operate on slash",
			options: Options{Args: []string{"/"}, TrashDir: trashDir},
			wantErr: true,
		},
		{
			name:    "remove empty dir with -d",
			options: Options{Args: []string{emptyDir}, TrashDir: trashDir, Dir: true},
			wantErr: false,
			check: func() {
				if _, err := os.Stat(filepath.Join(trashDir, "emptydir")); err != nil {
					t.Errorf("emptydir not moved to trash: %v", err)
				}
			},
		},
		{
			name:    "refuse to remove non-empty dir with -d",
			options: Options{Args: []string{nonEmptyDir}, TrashDir: trashDir, Dir: true},
			wantErr: true,
		},
		{
			name:    "remove non-empty dir recursively with -r",
			options: Options{Args: []string{nonEmptyDir}, TrashDir: trashDir, Recursive: true},
			wantErr: false,
			check: func() {
				if _, err := os.Stat(filepath.Join(trashDir, "nonemptydir")); err != nil {
					t.Errorf("nonemptydir not moved to trash: %v", err)
				}
			},
		},
		{
			name:    "non-existent file without force",
			options: Options{Args: []string{filepath.Join(tmpDir, "nope.txt")}, TrashDir: trashDir},
			wantErr: true,
		},
		{
			name:    "non-existent file with force",
			options: Options{Args: []string{filepath.Join(tmpDir, "nope.txt")}, TrashDir: trashDir, Force: true},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.options.Run()
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.check != nil {
				tt.check()
			}
		})
	}
}
