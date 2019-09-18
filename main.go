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

	tmpDir := fmt.Sprintf("%s/tmp", addonsDir)
	_ = os.Mkdir(tmpDir, os.ModePerm)

	curse := Curse(tmpDir, debug)
	unpacker := NewUnpacker()

	var wg sync.WaitGroup

	for _, addon := range cfg.CurseForge {
		wg.Add(1)

		go func(addon string) {
			defer wg.Done()
			file, sum, err := curse.Download(addon)
			must(err)

			unpacker.AddFile(addon, file, sum)
		}(addon)
	}

	gh := GitHub(tmpDir)

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
