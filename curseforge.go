package main

import "fmt"

const baseURL = "https://www.curseforge.com"

type CurseForgeDownloader struct {
	path     string
	debug    bool
	headless bool
}

func Curse(path string, debug bool, headless bool) *CurseForgeDownloader {
	return &CurseForgeDownloader{
		path:     path,
		debug:    debug,
		headless: headless,
	}
}

func (cfd *CurseForgeDownloader) DownloadFile(url string) (string, string, error) {
	return NewChrome(true).DownloadFileUsingChrome(url, cfd.path)
}

func FormatURL(name string) string {
	return fmt.Sprintf(`%s/wow/addons/%s/download`, baseURL, name)
}
