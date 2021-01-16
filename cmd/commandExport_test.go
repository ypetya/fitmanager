package main

import (
	"testing"

	"github.com/ypetya/fitmanager/cmd/mocks"
)

func Test_CommandExport_With_RemoteName_Should_Call_DataStore_Export(t *testing.T) {
	// GIVEN
	ds := mocks.FakeDataStore{LoadReturn: true}
	// WHEN
	commandExport([]string{"dir", "remote1"}, &ds)
	// THEN
	if ds.Calls[0] != "Load" {
		t.Error("Should call load first")
	}
	if ds.Calls[1] != "Export" || ds.Passed[0] != "remote1" {
		t.Error("Should call export with remote1!")
	}

	if ds.Calls[2] != "Save" {
		t.Error("Should call save after successful export")
	}
}
