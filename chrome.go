// not used file atm
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
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

func (c *Chrome) GetDownlaodHrefUsingChrome(url string) (string, error) {
	var href string
	var err error

	log.Printf("Navigating to %s", url)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	opts := c.chromeOpts()

	alloCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	logOpts := chromedp.WithErrorf(log.Printf)

	taskCtx, cancel := chromedp.NewContext(alloCtx, logOpts)
	defer cancel()

	err = chromedp.Run(taskCtx,
		page.SetDownloadBehavior(page.SetDownloadBehaviorBehaviorDeny),
		chromedp.Navigate(url),
	)
	if err != nil {
		return "", fmt.Errorf("Could not navigate to %s: %s", url, err)
	}

	for i := 0; i < 10; i++ {
		log.Println("Trying to fetch href")

		time.Sleep(time.Second * 5)

		err = chromedp.Run(taskCtx,
			chromedp.Evaluate(`$("a:contains('here')").attr('href')`, &href),
		)
		log.Printf(`Evaluate result "%s": "%s"`, href, err)

		if err != nil {
			continue
		}

		if href != "" {
			break
		}
	}

	if err != nil {
		return "", fmt.Errorf("Could get download url from %s: %s", url, err)
	}

	return href, nil

}
