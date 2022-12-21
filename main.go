package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	addonsDir string
	debug     bool
	fast      bool
	headless  bool
)

func init() {
	flag.StringVar(&addonsDir, "addons-dir", "./addons", "Addons directory")
	flag.BoolVar(&fast, "fast", false, "Run everything in parallel")
	flag.BoolVar(&debug, "debug", false, "Debug output")
	flag.BoolVar(&headless, "headless", false, "Run without Chrome GUI")
}

type downloadedAddon struct {
	addon string
	file  string
	sum   string
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	cfg, err := ParseConfig(fmt.Sprintf("%s/addons.yaml", addonsDir))

	must(err)

	tmpDir := fmt.Sprintf("%s/tmp", addonsDir)
	_ = os.Mkdir(tmpDir, os.ModePerm)

	curse := Curse(tmpDir, debug, headless)
	unpacker := NewUnpacker()

	var wg sync.WaitGroup

	urls := cfg.URLs

	for _, addon := range cfg.CurseForge {
		urls = append(urls, FormatURL(addon))
	}

	for _, url := range urls {
		wg.Add(1)

		f := func(url string) {
			defer wg.Done()
			file, sum, err := curse.DownloadFile(url)
			must(err)

			unpacker.AddFile(url, file, sum)
		}

		if fast {
			go f(url)
		} else {
			f(url)
		}
	}

	gh := GitHub(addonsDir)

	for folder, name := range cfg.GitHub {
		wg.Add(1)

		go func(folder, name string) {
			defer wg.Done()
			err := gh.Update(folder, name)
			must(err)
		}(folder, name)
	}

	wg.Wait()

	defer unpacker.Cleanup()

	err = unpacker.Unpack(addonsDir)
	must(err)

	err = unpacker.SaveCache(addonsDir)
	must(err)
}
