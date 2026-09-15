package main

import (
	"fmt"
	"io"
	"net/http"
)

const (
	KB = 1024
	MB = KB * 1024 // 1048576
	GB = MB * 1024 // 1073741824
)

type HTTPDownloader struct {
	httpClient  *http.Client
	fileURL     string
	fileHashURL string
	fileName    string
	maxFileSize int64
}

func NewHTTPDownLoader(
	httpClient *http.Client,
) *HTTPDownloader {
	return &HTTPDownloader{
		httpClient:  httpClient,
		maxFileSize: GB,
	}
}

func (hd *HTTPDownloader) CanDownload(url string) error {
	resp, err := hd.httpClient.Head(url)
	if err != nil {
		return fmt.Errorf("не удалось проверить URL: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("файл не найден или недоступен: %s", resp.Status)
	}
	if resp.ContentLength > hd.maxFileSize {
		return fmt.Errorf("file too large: %d bytes", resp.ContentLength)
	}
	return nil
}

func (hd *HTTPDownloader) Download(url string) ([]byte, error) {
	resp, err := hd.httpClient.Get(url)
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
