package main

import (
	"errors"
	"net/url"
	"testing"
)

func TestNewURLBuilder_PositiveCases(t *testing.T) {
	tests := []struct {
		name     string
		origin   string
		template *url.URL
	}{
		{
			"valid_http",
			"http://example.com",
			&url.URL{Scheme: "http", Host: "example.com"},
		},
		{
			"valid_https",
			"https://example.com",
			&url.URL{Scheme: "https", Host: "example.com"},
		},
		{
			"valid_with_port",
			"https://example.com:8080",
			&url.URL{Scheme: "https", Host: "example.com:8080"},
		},
		{
			"valid_http_with_trailing_slash",
			"http://example.com/",
			&url.URL{Scheme: "http", Host: "example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ans, err := NewURLBuilder(tt.origin)
			if err != nil {
				t.Errorf("NewURLBuilder(%q) error = %v", tt.origin, err)
			} else {
				assertURLsEqual(t, ans.urlTemplate, tt.template)
			}
		})
	}
}

func TestNewURLBuilder_NegativeCases(t *testing.T) {
	tests := []struct {
		name   string
		origin string
		err    error
	}{
		{
			"with_path",
			"http://example.com/my-path",
			ErrOriginContainsPath,
		},
		{
			"with_query",
			"https://example.com?a=1&b=2",
			ErrOriginContainsQuery,
		},
		{
			"valid_fragment",
			"https://example.com?#myfrag",
			ErrOriginContainsFragment,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ans, err := NewURLBuilder(tt.origin)
			if !errors.Is(err, tt.err) {
				t.Errorf("NewURLBuilder(%q) expected error", tt.origin)
				return
			}
			if ans != nil {
				t.Errorf("NewURLBuilder(%q) expected empty instance", tt.origin)
				return
			}
		})
	}
}

func TestURLBuilder_Build(t *testing.T) {
	tests := []struct {
		name     string
		origin   string
		parts    []string
		wantPath string
	}{
		{"no parts", "https://example.com", nil, "/"},
		{"one part", "https://example.com", []string{"users"}, "/users"},
		{"two parts", "https://example.com", []string{"users", "123"}, "/users/123"},
		{
			"with base path",
			"https://example.com/api",
			[]string{"v1", "users"},
			"/api/v1/users",
		},
		{
			"with special chars",
			"https://example.com",
			[]string{"users", "john doe"},
			"/users/john doe",
		},
		{
			"with slash in part",
			"https://example.com",
			[]string{"users", "a/b"},
			"/users/a/b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ub, err := NewURLBuilder(tt.origin)
			if err != nil {
				t.Fatalf("failed to create builder: %v", err)
			}

			got := ub.Build(tt.parts...)
			if got.Path != tt.wantPath {
				t.Errorf("Build() path = %q, want %q", got.Path, tt.wantPath)
			}
		})
	}
}

func TestURLBuilder_AsString(t *testing.T) {
	tests := []struct {
		name   string
		origin string
		parts  []string
		want   string
	}{
		{"no parts", "https://example.com", nil, "https://example.com"},
		{
			"one part",
			"https://example.com",
			[]string{"users"},
			"https://example.com/users",
		},
		{
			"two parts",
			"https://example.com",
			[]string{"users", "123"},
			"https://example.com/users/123",
		},
		{
			"with base path",
			"https://example.com/api",
			[]string{"v1", "users"},
			"https://example.com/api/v1/users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ub, err := NewURLBuilder(tt.origin)
			if err != nil {
				t.Fatalf("failed to create builder: %v", err)
			}

			got := ub.AsString(tt.parts...)
			if got != tt.want {
				t.Errorf("AsString() = %q, want %q", got, tt.want)
			}
		})
	}
}
