package mocks

import (
	"io/ioutil"
	"os"
)

var tempFiles []*os.File

type FakeIO struct {
	LastCall         string
	LastCallWithName string
	LastCallDst      string
	F                *os.File
}

func (f *FakeIO) CreateFileForWriting(name string) *os.File {
	f.LastCall, f.LastCallWithName = "createFile", name
	t, err := ioutil.TempFile("testdata", "test")
	bail(err)
	tempFiles = append(tempFiles, t)
	f.F = t
	return t
}

func (f *FakeIO) Copy(src, dst string) {
	f.LastCall, f.LastCallWithName, f.LastCallDst = "Copy", src, dst
}

// This is a real implementation as we do not mock read operations
func (f *FakeIO) FileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func (f *FakeIO) SameSize(fileName string) bool {
	f1, err := os.Stat(f.F.Name())
	bail(err)
	f2, err := os.Stat(fileName)
	bail(err)
	return f1.Size() == f2.Size()
}

func FakeIODeleteAllTempfiles() {
	for _, val := range tempFiles {
		os.Remove(val.Name())
	}
}
