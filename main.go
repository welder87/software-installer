package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	KB = 1024
	MB = KB * 1024 // 1048576
	GB = MB * 1024 // 1073741824
)
const latestTagTemplate = "https://api.github.com/repos/%s/%s/releases/latest"

type HTTPLoader struct {
	httpClient  *http.Client
	fileURL     string
	fileHashURL string
	fileName    string
	maxFileSize int64
}

func NewHTTPLoader(
	httpClient *http.Client,
	fileURL string,
	fileHashURL string,
	fileName string,
) *HTTPLoader {
	return &HTTPLoader{
		httpClient:  httpClient,
		fileURL:     fileURL,
		fileHashURL: fileHashURL,
		fileName:    fileName,
		maxFileSize: GB,
	}
}

func (hl *HTTPLoader) extractHash(reader io.ReadCloser) (string, error) {
	defer reader.Close()
	scanner := bufio.NewScanner(reader)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		if strings.Contains(line, hl.fileName) {
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

func canDownloadFactory(httpClient *http.Client, maxFileSize int64) func(string) error {
	return func(url string) error {
		resp, err := httpClient.Head(url)
		if err != nil {
			return fmt.Errorf("не удалось проверить URL: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("файл не найден или недоступен: %s", resp.Status)
		}
		if resp.ContentLength > maxFileSize {
			return fmt.Errorf("file too large: %d bytes", resp.ContentLength)
		}
		return nil
	}
}

func downloadFactory(httpClient *http.Client) func(string) ([]byte, error) {
	return func(url string) ([]byte, error) {
		resp, err := httpClient.Get(url)
		if err != nil {
			return nil, fmt.Errorf("не удалось скачать файл: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("status: %s", resp.Status)
		}
		result, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("ошибка при чтении тела ответа: %w", err)
		}
		return result, nil
	}
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

/*
LAZYGIT_VERSION=$(curl -s "https://api.github.com/repos/jesseduffield/lazygit/releases/latest" | \grep -Po '"tag_name": *"v\K[^"]*')
LAZYGIT_ARCH=$(uname -m | sed -e 's/aarch64/arm64/')
curl -Lo lazygit.tar.gz "https://github.com/jesseduffield/lazygit/releases/download/v${LAZYGIT_VERSION}/lazygit_${LAZYGIT_VERSION}_Linux_${LAZYGIT_ARCH}.tar.gz"
tar xf lazygit.tar.gz lazygit
sudo install lazygit -D -t /usr/local/bin/
*/
func main() {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		DisableKeepAlives:   false, // true = отключить keep-alive для тестов.
	}
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
	canDownload := canDownloadFactory(httpClient, 100*MB)
	download := downloadFactory(httpClient)
	softData := [][]string{
		{"zellij-org", "zellij"},
		{"charmbracelet", "glow"},
	}
	for _, items := range softData {
		latestReleaseURL := fmt.Sprintf(latestTagTemplate, items[0], items[1])
		err := canDownload(latestReleaseURL)
		if err != nil {
			fmt.Println(err)
			continue
		}
		res, err := download(latestReleaseURL)
		if err != nil {
			fmt.Println(err)
			continue
		}
		nres, err := extractLatestTag(res)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(nres)
	}
}
