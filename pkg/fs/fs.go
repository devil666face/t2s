package fs

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
)

func WriteFile(
	path string,
	data []byte,
	perm ...int,
) error {
	if len(perm) == 0 {
		perm = append(perm, 0644)
	}
	if err := os.WriteFile(path, data, fs.FileMode(perm[0])); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}
	return nil
}

func ReadFile(path string) ([]byte, error) {
	f, _, err := OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var buf bytes.Buffer
	if _, err = buf.ReadFrom(f); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func OpenFile(filename string) (*os.File, fs.FileInfo, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("error file %s not open: %w", filename, err)
	}
	stat, err := f.Stat()
	if err != nil {
		return nil, nil, fmt.Errorf("error to get file size: %w", err)
	}
	return f, stat, nil
}

func CheckSymlink(path string) (string, error) {
	target, err := os.Readlink(path)
	if err != nil {
		return "", err
	}
	return target, nil
}
