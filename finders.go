package main

import (
	"errors"
	"fmt"
)

var ErrAssetNotFound = errors.New("asset not found")

func findProgram(
	release LatestReleaseInfo,
	filename string,
) (LatestReleaseAsset, error) {
	for _, asset := range release.Assets {
		if string(asset.Name) == filename {
			return asset, nil
		}
	}
	return LatestReleaseAsset{}, fmt.Errorf(
		"asset with filename %s not found %w",
		filename,
		ErrAssetNotFound,
	)
}
