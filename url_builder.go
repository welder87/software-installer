package main

import (
	"errors"
	"fmt"
	"net/url"
)

var (
	ErrOriginContainsPath     = errors.New("origin must not contain a path")
	ErrOriginContainsQuery    = errors.New("origin must not contain a query")
	ErrOriginContainsFragment = errors.New("origin must not contain a fragment")
	ErrEmptyPath              = errors.New("path parts must not be empty")
)

type URLBuilder struct {
	urlTemplate *url.URL
}

func NewURLBuilder(origin string) (*URLBuilder, error) {
	item, err := url.Parse(origin)
	if err != nil {
		return nil, fmt.Errorf("invalid origin: %s. %w", origin, err)
	}
	if item.Path != "" && item.Path != "/" {
		return nil, fmt.Errorf("invalid origin: %s. %w", origin, ErrOriginContainsPath)
	}
	if item.RawQuery != "" {
		return nil, fmt.Errorf("invalid origin: %s. %w", origin, ErrOriginContainsQuery)
	}
	if item.Fragment != "" {
		return nil, fmt.Errorf(
			"invalid origin: %s. %w",
			origin,
			ErrOriginContainsFragment,
		)
	}

	item.Path = ""

	return &URLBuilder{
		urlTemplate: item,
	}, nil
}

func (ub *URLBuilder) Build(item ...string) (*url.URL, error) {
	if len(item) == 0 {
		return nil, ErrEmptyPath
	}
	return ub.urlTemplate.JoinPath(item...), nil
}

func (ub *URLBuilder) AsString(item ...string) (string, error) {
	res, err := ub.Build(item...)
	if err != nil {
		return "", err
	}
	return res.String(), nil
}
