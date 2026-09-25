package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
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

func calculateHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open file %s: %w", path, err)
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("stat file %s: %w", path, err)
	}
	if fi.IsDir() {
		return "", fmt.Errorf("%s is a directory, not a file", path)
	}
	h := sha256.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return "", fmt.Errorf("read file %s for hashing: %w", path, err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
