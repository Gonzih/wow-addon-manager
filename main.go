package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
)

var addonsDir string
var debug bool

func init() {
	flag.StringVar(&addonsDir, "addons-dir", "./addons", "Addons directory")
	flag.BoolVar(&debug, "debug", false, "Debug output")
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
	log.Printf("%#v", cfg)

	tmpDir := fmt.Sprintf("%s/tmp", addonsDir)
	_ = os.Mkdir(tmpDir, os.ModePerm)

	downloader := Curse(tmpDir, debug)
	unpacker := NewUnpacker()

	var wg sync.WaitGroup

	for _, addon := range cfg.Addons {
		wg.Add(1)

		go func(addon string) {
			defer wg.Done()
			file, sum, err := downloader.Download(addon)
			must(err)

			unpacker.AddFile(addon, file, sum)
		}(addon)
	}

	wg.Wait()

	defer unpacker.Cleanup()

	err = unpacker.Unpack(addonsDir)
	must(err)

	err = unpacker.SaveCache(addonsDir)
	must(err)
}
