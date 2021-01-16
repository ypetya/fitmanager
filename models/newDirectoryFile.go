package models

import (
	"os"
	"time"
)

// Checks the directory existance, if it exists, fills out the lastSeen field
func newFromDir(path string) File {
	var lastSeen int64
	if directoryExists(path) {
		lastSeen = time.Now().Unix()
	}
	return File{Path: path, LastSeen: lastSeen}
}

func directoryExists(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		// file not exist
		return false
	}
	_, err = f.Readdir(-1)
	f.Close()
	if err != nil {
		// dir not exist
		return false
	}

	return true
}
