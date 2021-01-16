package internal

import (
	"io"
	"os"
)

// A concrete implementation for file operations
// Implements IFileCreator
type EffectiveIO struct {
}

func (e EffectiveIO) CreateFileForWriting(name string) *os.File {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	bail(err)
	return f
}

func (e EffectiveIO) FileExists(name string) bool {
	_, err := os.Stat(name)
	// This is not te same as os.IsExist !!! as it checks the error !!!
	return !os.IsNotExist(err)
}

func (e EffectiveIO) Copy(src, dst string) {
	sf, err := os.Open(src)
	defer sf.Close()
	bail(err)
	tf := e.CreateFileForWriting(dst)
	defer tf.Close()
	_, err = io.Copy(tf, sf)
	bail(err)
}
