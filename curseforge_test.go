package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// func TestDownloadUrl(t *testing.T) {
// 	curse := Curse("./addons/tmp", false)
// 	url, err := curse.getDownloadUrl("atlas")
// 	require.Nil(t, err)
// 	require.Regexp(t, regexp.MustCompile("^https://www.curseforge.com/wow/addons/atlas/download/\\d+/file$"), url)
// }

func TestDownloadFile(t *testing.T) {
	url := fmt.Sprintf(`%s/wow/addons/%s/download`, baseURL, "atlas")
	curse := Curse("./addons/tmp", false)
	path, _, err := curse.DownloadFile(url)
	defer os.Remove(path)
	require.Nil(t, err)
	require.FileExists(t, path)
}
