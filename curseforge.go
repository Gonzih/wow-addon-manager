package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/antchfx/htmlquery"
	"github.com/google/uuid"
)

const baseURL = "https://www.curseforge.com"

type CurseForgeDownloader struct {
	path  string
	debug bool
}

func Curse(path string, debug bool) *CurseForgeDownloader {
	return &CurseForgeDownloader{
		path:  path,
		debug: debug,
	}
}

func (cfd *CurseForgeDownloader) getDownloadUrl(name string) (string, error) {
	url := fmt.Sprintf(`%s/wow/addons/%s/download`, baseURL, name)
	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		return "", fmt.Errorf("Could get download url for %s: %s", name, err)
	}

	xpath := `//a[text()='here']`
	a := htmlquery.FindOne(doc, xpath)
	href := htmlquery.SelectAttr(a, "href")

	// if href == "" {
	// 	log.Printf("Trying to use chrome to download url since href was empty")
	// 	chrome := NewChrome(true)
	// 	href, err = chrome.GetDownlaodHrefUsingChrome(url, xpath)
	// }

	href = fmt.Sprintf("%s%s", baseURL, href)

	return href, nil
}

func (cfd *CurseForgeDownloader) downloadFile(url string) (string, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	fname := uuid.New().String()
	filepath := fmt.Sprintf("%s/%s.zip", cfd.path, fname)

	out, err := os.Create(filepath)
	if err != nil {
		return "", "", err
	}
	defer out.Close()

	mdbuf := md5.New()
	_, err = io.Copy(out, io.TeeReader(resp.Body, mdbuf))
	if err != nil {
		return "", "", err
	}

	sum := fmt.Sprintf("%x", mdbuf.Sum(nil))

	return filepath, sum, nil
}

func (cfd *CurseForgeDownloader) Download(name string) (string, string, error) {
	url, err := cfd.getDownloadUrl(name)
	if err != nil {
		return "", "", err
	}

	log.Printf("Going to url %s to download %s", url, name)

	return cfd.downloadFile(url)
}
