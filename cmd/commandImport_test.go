package main

import (
	"testing"

	"github.com/ypetya/fitmanager/cmd/mocks"
)

func TestCommandImportWithoutRemoteNameShouldImportFromGarmin(t *testing.T) {
	// GIVEN
	ds := mocks.FakeDataStore{LoadReturn: true}
	// WHEN
	commandImport([]string{"dir"}, &ds)
	// THEN
	if ds.Calls[0] != "Load" {
		t.Error("Should call load first")
	}
	if ds.Calls[1] != "Import" || ds.Passed[0] != "garmin" {
		t.Error("Should call import with remote garmin by default!")
	}

	if ds.Calls[2] != "Save" {
		t.Error("Should call save after successful import")
	}
}
