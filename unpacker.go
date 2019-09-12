package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v2"
)

type AddonsCache struct {
	Checksums map[string]string `yaml:"checksums"`
}

type FileInfo struct {
	Addon string
	File  string
	Sum   string
}

type Unpacker struct {
	Files     []FileInfo
	filesLock sync.Mutex
	Cache     *AddonsCache
}

func NewUnpacker() *Unpacker {
	return &Unpacker{}
}

func (u *Unpacker) AddFile(addon, file, sum string) {
	u.filesLock.Lock()
	defer u.filesLock.Unlock()

	u.Files = append(u.Files, FileInfo{addon, file, sum})
}

func (u *Unpacker) Unpack(dest string) error {
	err := u.LoadCache(dest)
	if err != nil {
		return fmt.Errorf("Error loading cache: %s", err)
	}

	for _, fInfo := range u.Files {
		filePath := fInfo.File
		addon := fInfo.Addon
		sum := fInfo.Sum

		cacheSum, ok := u.Cache.Checksums[addon]

		if ok && sum == cacheSum {
			continue
		}

		log.Printf("Updating %s", addon)
		err = Unzip(filePath, dest)
		if err != nil {
			return fmt.Errorf("Error unpacking %s: %s", filePath, err)
		}

		u.Cache.Checksums[addon] = sum
	}

	err = u.SaveCache(dest)
	if err != nil {
		return fmt.Errorf("Error saving cache: %s", err)
	}

	return nil
}
func (u *Unpacker) cachePath(dest string) string {
	return filepath.Join(dest, "cache.yaml")
}

func (u *Unpacker) LoadCache(dest string) error {
	cachePath := u.cachePath(dest)

	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		u.Cache = &AddonsCache{Checksums: make(map[string]string, 0)}
		return nil
	}

	file, err := os.Open(cachePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fdata, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	cache := &AddonsCache{}
	err = yaml.Unmarshal(fdata, cache)
	if err != nil {
		return err
	}

	u.Cache = cache

	return nil
}

func (u *Unpacker) SaveCache(dest string) error {
	cachePath := u.cachePath(dest)

	file, err := os.Create(cachePath)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(u.Cache)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (u *Unpacker) Cleanup() {
	for _, fInfo := range u.Files {
		err := os.Remove(fInfo.File)
		if err != nil {
			log.Printf("Could not delete file for %s: %s", fInfo.Addon, err)
		}
	}
}
