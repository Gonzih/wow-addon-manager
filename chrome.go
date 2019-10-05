// not used file atm
package main

import (
	"context"
	"fmt"
	"log"
	"time"

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

func (c *Chrome) GetDownlaodHrefUsingChrome(url, xpath string) (string, error) {
	log.Printf("Navigating to %s", url)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := c.chromeOpts()

	alloCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	logOpts := chromedp.WithErrorf(log.Printf)

	taskCtx, cancel := chromedp.NewContext(alloCtx, logOpts)
	defer cancel()

	var href string
	var ok bool
	err := chromedp.Run(taskCtx,
		chromedp.Navigate(url),
		chromedp.AttributeValue(xpath, "href", &href, &ok, chromedp.BySearch),
	)

	if err != nil {
		return "", fmt.Errorf("Could get download url from %s: %s", url, err)
	}

	if !ok {
		return "", fmt.Errorf("Could not get download url: %s", href)
	}

	return href, nil

}
