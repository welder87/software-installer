package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

func extractLatestReleaseInfo(lines []byte) (LatestReleaseInfo, error) {
	var release LatestReleaseInfo
	err := json.Unmarshal(lines, &release)
	if err != nil {
		return LatestReleaseInfo{}, fmt.Errorf("contract violation: %w", err)
	}
	return release, nil
}

type LatestReleaseInfo struct {
	TagName TrimmedString         `json:"tag_name"`
	Assets  []LatestReleaseAsset `json:"assets"`
}

type LatestReleaseAsset struct {
	ContentType        TrimmedString  `json:"content_type"`
	State              TrimmedString  `json:"state"`
	Name               TrimmedString  `json:"name"`
	Digest             TrimmedString  `json:"digest"`
	BrowserDownloadURL MarshalableURL `json:"browser_download_url"`
}

type MarshalableURL struct {
	*url.URL
}

func (u *MarshalableURL) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "null" {
		u.URL = nil
		return nil
	}
	parsed, err := url.Parse(s)
	if err != nil {
		return err
	}
	u.URL = parsed
	return nil
}

func (u MarshalableURL) MarshalJSON() ([]byte, error) {
	if u.URL == nil {
		return []byte("null"), nil
	}
	return json.Marshal(u.URL.String())
}

type TrimmedString string

func (s *TrimmedString) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*s = TrimmedString(strings.TrimSpace(raw))
	return nil
}

func (s TrimmedString) MarshalJSON() ([]byte, error) {
	return json.Marshal(strings.TrimSpace(string(s)))
}
