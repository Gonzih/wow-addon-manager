package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var addonsDir string
var debug bool

func init() {
	flag.StringVar(&addonsDir, "addons-dir", "./addons", "Addons directory")
	flag.BoolVar(&debug, "debug", false, "Debug output")
	flag.Parse()
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
	cfg, err := ParseConfig(fmt.Sprintf("%s/addons.yaml", addonsDir))

	must(err)
	log.Printf("%#v", cfg)

	tmpDir := fmt.Sprintf("%s/tmp", addonsDir)
	_ = os.Mkdir(tmpDir, os.ModePerm)

	downloader := Curse(tmpDir, debug)
	unpacker := NewUnpacker()

	for _, addon := range cfg.Addons {
		file, sum, err := downloader.Download(addon)
		must(err)

		unpacker.AddFile(addon, file, sum)
	}

	err = unpacker.Unpack()
	must(err)

	log.Printf("%#v", unpacker)
}
