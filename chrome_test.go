package main

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDownloadHrefWithChrome(t *testing.T) {
	url := fmt.Sprintf(`%s/wow/addons/%s/download`, baseURL, "atlas")

	chrome := NewChrome(true)
	href, err := chrome.GetDownlaodHrefUsingChrome(url)

	require.Nil(t, err)
	require.Regexp(t, regexp.MustCompile("^/wow/addons/atlas/download/\\d+/file$"), href)
}
