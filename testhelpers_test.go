package main

import (
	"net/url"
	"testing"
)

func assertURLsEqual(t *testing.T, got, want *url.URL) {
	t.Helper()

	equal, diffs := compareURLs(got, want)
	if !equal {
		t.Errorf(
			"URL mismatch:\n  got:  %s\n  want: %s\n  diffs: %v",
			got.String(),
			want.String(),
			diffs,
		)
	}
}

func compareURLs(u1, u2 *url.URL) (bool, []string) {
	var diffs []string
	if u1.Scheme != u2.Scheme {
		diffs = append(diffs, "scheme")
	}
	if u1.Host != u2.Host {
		diffs = append(diffs, "host")
	}
	if u1.Path != u2.Path {
		diffs = append(diffs, "path")
	}
	if u1.Query().Encode() != u2.Query().Encode() {
		diffs = append(diffs, "query")
	}
	return len(diffs) == 0, diffs
}
