package main

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestComposePageUrl(t *testing.T) {
	baseUrl := url.URL{
		Scheme: "https",
		Host:   "example.com",
		Path:   "/people"}

	urlStr := composePageUrl(baseUrl, 3, 20)

	assert.Equal(t, urlStr, "https://example.com/people?limit=20&page=3")

	baseUrl = url.URL{
		Scheme:   "http",
		Host:     "example.com",
		Path:     "/relations",
		RawQuery: "flag=t"}

	urlStr = composePageUrl(baseUrl, 3, 20)

	assert.Equal(t, urlStr, "http://example.com/relations?flag=t&limit=20&page=3")
}
