package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
)

const baseURL = "https://www.curseforge.com"

type CurseForgeDownloader struct {
	path     string
	debug    bool
	headless bool
}

func Curse(path string, debug, headless bool) *CurseForgeDownloader {
	return &CurseForgeDownloader{
		path:     path,
		debug:    debug,
		headless: headless,
	}
}
func (cfd *CurseForgeDownloader) chromeOpts() []func(*chromedp.ExecAllocator) {
	opts := []func(*chromedp.ExecAllocator){
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,

		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("enable-features", "NetworkService,NetworkServiceInProcess"),
		chromedp.Flag("disable-background-timer-throttling", true),
		chromedp.Flag("disable-backgrounding-occluded-windows", true),
		chromedp.Flag("disable-breakpad", true),
		chromedp.Flag("disable-client-side-phishing-detection", true),
		chromedp.Flag("disable-default-apps", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-features", "site-per-process,TranslateUI,BlinkGenPropertyTrees"),
		chromedp.Flag("disable-hang-monitor", true),
		chromedp.Flag("disable-ipc-flooding-protection", true),
		chromedp.Flag("disable-popup-blocking", true),
		chromedp.Flag("disable-prompt-on-repost", true),
		chromedp.Flag("disable-renderer-backgrounding", true),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("force-color-profile", "srgb"),
		chromedp.Flag("metrics-recording-only", true),
		chromedp.Flag("safebrowsing-disable-auto-update", true),
		chromedp.Flag("enable-automation", true),
		chromedp.Flag("password-store", "basic"),
		chromedp.Flag("use-mock-keychain", true),
		chromedp.DisableGPU,
	}

	if cfd.headless {
		opts = append(opts,
			chromedp.Headless,
		)
	}

	return opts
}

func (cfd *CurseForgeDownloader) getDownloadUrl(name string) (string, error) {

	log.Printf("Trying to download %s", name)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	opts := cfd.chromeOpts()

	alloCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	logOpts := chromedp.WithErrorf(log.Printf)

	if cfd.debug {
		// logOpts = chromedp.WithDebugf(log.Printf)
	}

	taskCtx, cancel := chromedp.NewContext(alloCtx, logOpts)
	defer cancel()

	var href string
	var ok bool
	err := chromedp.Run(taskCtx,
		chromedp.Navigate(fmt.Sprintf(`%s/wow/addons/%s/download`, baseURL, name)),
		chromedp.AttributeValue(`//a[text()='here']`, "href", &href, &ok),
	)

	if err != nil {
		return "", fmt.Errorf("Could get download url for %s: %s", name, err)
	}

	if !ok {
		return "", fmt.Errorf("Could not get download url: %s", href)
	}

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

	cfd.Logf("Going to url %s to download %s", url, name)

	return cfd.downloadFile(url)
}

func (cfd *CurseForgeDownloader) Logf(s string, args ...interface{}) {
	if cfd.debug {
		log.Printf(s, args...)
	}
}
