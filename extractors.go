package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func extractHash(reader io.ReadCloser, filename string) (string, error) {
	defer reader.Close()
	scanner := bufio.NewScanner(reader)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		if strings.Contains(line, filename) {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	parts := strings.SplitN(line, " ", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("неверный формат строки")
	}
	return parts[1], nil
}

type latestReleaseInfo struct {
	TagName string `json:"tag_name"`
}

func extractLatestTag(lines []byte) (latestReleaseInfo, error) {
	var release latestReleaseInfo
	err := json.Unmarshal(lines, &release)
	if err != nil {
		return latestReleaseInfo{}, fmt.Errorf("ошибка при парсинге JSON: %w", err)
	}
	return release, nil
}
