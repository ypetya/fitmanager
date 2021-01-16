package internal

import "os"

type IFileCreator interface {
	CreateFileForWriting(string) *os.File
	FileExists(string) bool
	Copy(src, dst string)
}
