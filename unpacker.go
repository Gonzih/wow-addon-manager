package main

type FileInfo struct {
	Addon string
	File  string
	Sum   string
}

type Unpacker struct {
	Files []FileInfo
}

func NewUnpacker() *Unpacker {
	return &Unpacker{}
}

func (u *Unpacker) AddFile(addon, file, sum string) {
	u.Files = append(u.Files, FileInfo{addon, file, sum})
}

func (u *Unpacker) Unpack() error {
	return nil
}
