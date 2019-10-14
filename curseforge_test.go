package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDownloadFile(t *testing.T) {
	url := fmt.Sprintf(`%s/wow/addons/%s/download`, baseURL, "atlas")
	curse := Curse("./addons/tmp", false)
	path, _, err := curse.DownloadFile(url)
	defer os.Remove(path)
	require.Nil(t, err)
	require.FileExists(t, path)
}
