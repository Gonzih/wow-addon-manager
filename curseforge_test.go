package main

import (
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDownloadUrl(t *testing.T) {
	curse := Curse("/tmp", false)
	url, err := curse.getDownloadUrl("bartender4")
	require.Nil(t, err)
	require.Regexp(t, regexp.MustCompile("^https://www.curseforge.com/wow/addons/bartender4/download/\\d+/file$"), url)
}

func TestDownloadFile(t *testing.T) {
	curse := Curse("./addons/tmp", false)
	path, _, err := curse.downloadFile("https://gonzih.me/index.html")
	defer os.Remove(path)
	require.Nil(t, err)
	require.FileExists(t, path)
}
