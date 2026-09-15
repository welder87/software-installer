package main

import (
	"fmt"
	"net/http"
	"time"
)

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
	httpDownloader := NewHTTPDownLoader(httpClient)
	softData := [][]string{
		{"zellij-org", "zellij"},
		{"charmbracelet", "glow"},
	}
	for _, items := range softData {
		urlBuilder, err := NewURLBuilder("https://api.github.com")
		if err != nil {
			fmt.Println("Error %w", err)
			continue
		}
		latestReleaseURL, err := urlBuilder.AsString(
			"repos",
			items[0],
			items[1],
			"releases",
			"latest",
		)
		if err != nil {
			fmt.Println("Error %w", err)
			continue
		}
		err = httpDownloader.CanDownload(latestReleaseURL)
		if err != nil {
			fmt.Println(err)
			continue
		}
		res, err := httpDownloader.Download(latestReleaseURL)
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
