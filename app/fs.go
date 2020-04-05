package app

import "github.com/gobuffalo/packr"

type Fs struct {
	Templates packr.Box
}

var (
	fs *Fs
)

func init() {
	fs = new(Fs)
}

func GetFS() *Fs {
	return fs
}
