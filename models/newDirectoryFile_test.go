package models

import (
	"testing"
)

func TestNewShouldHaveLastSeenWhenDirExist(t *testing.T) {
	f := newFromDir("testdata/empty")
	if f.LastSeen == 0 {
		t.Error("Expected File.LastSeen to set")
	}
}
func TestNewShouldNotSetLastSeenWhenNoSuchFileExist(t *testing.T) {
	f := newFromDir("testdata/undefined")
	if f.LastSeen != 0 {
		t.Error("Expected File.LastSeen to NOT set")
	}
}
func TestNewShouldNotSetLastSeenWhenItIsNotADir(t *testing.T) {
	f := newFromDir("testdata/file.txt")
	if f.LastSeen != 0 {
		t.Error("Expected File.LastSeen to NOT set")
	}
}
