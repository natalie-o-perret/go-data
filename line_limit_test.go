package data_test

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestGoFilesDoNotExceed100Lines(t *testing.T) {
	const limit = 100
	err := filepath.WalkDir(".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if entry.Name() == ".git" || entry.Name() == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}

		contents, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		lines := bytes.Count(contents, []byte{'\n'})
		if len(contents) > 0 && contents[len(contents)-1] != '\n' {
			lines++
		}
		if lines > limit {
			t.Errorf("%s has %d physical lines, maximum is %d",
				filepath.ToSlash(path), lines, limit)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
