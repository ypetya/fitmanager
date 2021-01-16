package main

import (
	"testing"

	"./mocks"
)

func TestCommandInit(t *testing.T) {
	// GIVEN
	ds := mocks.FakeDataStore{
		LoadReturn: false,
	}
	// WHEN
	commandInit([]string{"dir"}, &ds)
	// THEN
	if ds.Calls[0] != "Load" {
		t.Error("Should call load first")
	}
	if ds.Calls[1] != "AddRemote" || ds.Passed[0] != "garmin" {
		t.Error("Should call add garmin as remote by default")
	}
	if ds.Calls[2] != "Save" {
		t.Error("Should call Save")
	}
}
