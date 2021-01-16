package internal

import "os"

type FakeIO struct {
	Called string
	Name   string
	Dst    string
}

func (f *FakeIO) CreateFileForWriting(name string) *os.File {
	f.Called, f.Name = "createFile", name
	return &os.File{}
}

// This is a real implementation as we do not mock read operations
func (f *FakeIO) FileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func (f *FakeIO) Copy(src, dst string) {
	f.Called, f.Name, f.Dst = "Copy", src, dst
}
