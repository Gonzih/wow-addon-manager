package main

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDownloadHrefWithChrome(t *testing.T) {
	url := fmt.Sprintf(`%s/wow/addons/%s/download`, baseURL, "bartender4")
	xpath := `//a[text()='here']`

	chrome := NewChrome(true)
	href, err := chrome.GetDownlaodHrefUsingChrome(url, xpath)

	require.Nil(t, err)
	require.Regexp(t, regexp.MustCompile("^/wow/addons/bartender4/download/\\d+/file$"), href)
}
