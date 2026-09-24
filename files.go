package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func makeFilePath(path string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("home dir error %w", err)
	}
	absPath := filepath.Join(home, path)
	isExists, err := fileExists(absPath)
	if err != nil {
		return "", fmt.Errorf("file error %w", err)
	}
	if !isExists {
		return "", fmt.Errorf("file not found")
	}
	return absPath, nil
}
