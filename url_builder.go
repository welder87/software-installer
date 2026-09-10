package main

import (
	"fmt"
	"net/url"
)

type URLBuilder struct {
	urlTemplate *url.URL
}

func NewURLBuilder(origin string) (*URLBuilder, error) {
	item, err := url.Parse(origin)
	if err != nil {
		return nil, fmt.Errorf("invalid origin %w", err)
	}
	if item.Path != "" && item.Path != "/" {
		return nil, fmt.Errorf("origin must not contain a path: %s", item.Path)
	}
	if item.RawQuery != "" {
		return nil, fmt.Errorf("origin must not contain a query: %s", item.RawQuery)
	}
	if item.Fragment != "" {
		return nil, fmt.Errorf("origin must not contain a fragment: %s", item.Fragment)
	}

	item.Path = ""

	return &URLBuilder{
		urlTemplate: item,
	}, nil
}

func (ub *URLBuilder) Build(item ...string) *url.URL {
	return ub.urlTemplate.JoinPath(item...)
}

func (ub *URLBuilder) AsString(item ...string) string {
	return ub.Build(item...).String()
}
