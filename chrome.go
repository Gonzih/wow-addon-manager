// not used file atm
package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
)

type Chrome struct {
	headless bool
}

func NewChrome(headless bool) *Chrome {
	return &Chrome{headless: headless}
}

func (c *Chrome) chromeOpts() []func(*chromedp.ExecAllocator) {
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

	if c.headless {
		opts = append(opts,
			chromedp.Headless,
		)
	}

	return opts
}

// func (c *Chrome) GetDownlaodHrefUsingChrome(url string) (string, error) {
// 	var href string
// 	var err error

// 	for i := 0; i < 10; i++ {

// 		log.Printf("Navigating to %s", url)

// 		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
// 		defer cancel()

// 		opts := c.chromeOpts()

// 		alloCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
// 		defer cancel()

// 		logOpts := chromedp.WithErrorf(log.Printf)

// 		taskCtx, cancel := chromedp.NewContext(alloCtx, logOpts)
// 		defer cancel()

// 		err = chromedp.Run(taskCtx,
// 			page.SetDownloadBehavior(page.SetDownloadBehaviorBehaviorDeny),
// 			chromedp.Navigate(url),
// 		)
// 		if err != nil {
// 			log.Printf("Could not navigate to %s: %s", url, err)
// 			continue
// 		}

// 		time.Sleep(time.Second * 15)

// 		log.Println("Trying to query for href")
// 		err = chromedp.Run(taskCtx,
// 			chromedp.Evaluate(`$('a').map((i, el) => $(el).attr('href')).toArray().find(h => h.match(/\/wow\/addons\/.+\/download\/\d+\/file/))`, &href),
// 		)
// 		log.Printf(`Evaluate result is "%s": "%v"`, href, err)

// 		if err != nil {
// 			continue
// 		}

// 		if href != "" {
// 			break
// 		}
// 	}

// 	if err != nil {
// 		return "", fmt.Errorf("Could get download url from %s: %s", url, err)
// 	}

// 	return href, nil

// }

func (c *Chrome) DownloadFileUsingChrome(url, tmp string) (string, string, error) {
	folderName := uuid.New().String()
	folder := fmt.Sprintf("%s/%s", tmp, folderName)
	os.MkdirAll(folder, os.ModePerm)

	log.Printf("Navigating to %s", url)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	opts := c.chromeOpts()

	alloCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	logOpts := chromedp.WithErrorf(log.Printf)

	taskCtx, cancel := chromedp.NewContext(alloCtx, logOpts)
	defer cancel()

	log.Printf("Setting download path to %s", folder)

	err := chromedp.Run(taskCtx,
		page.SetDownloadBehavior(page.SetDownloadBehaviorBehaviorAllow).WithDownloadPath(folder),
		chromedp.Navigate(url),
	)

	time.Sleep(time.Second * 20)

	if err != nil {
		log.Printf(`Chrome session exit error that will be ignored: "%v"`, err)
	}

	files, err := filepath.Glob(fmt.Sprintf("%s/*", folder))
	if err != nil {
		return "", "", fmt.Errorf("Could not find files in %s: %s", folder, err)
	}

	if len(files) == 0 {
		return "", "", fmt.Errorf("Folder %s is empty", folder)
	}

	var file string
	reg := regexp.MustCompile(`\.zip$`)
	for _, f := range files {
		if reg.MatchString(f) {
			file = f
			break
		}
	}

	mdbuf := md5.New()
	f, err := os.Open(file)
	if err != nil {
		return "", "", fmt.Errorf("Could not open file %s: %s", file, err)
	}

	_, err = io.Copy(mdbuf, f)
	if err != nil {
		return "", "", fmt.Errorf("Could not read file %s to calculate checksum: %s", file, err)
	}

	sum := fmt.Sprintf("%x", mdbuf.Sum(nil))

	return file, sum, nil
}
