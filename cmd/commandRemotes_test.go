package main

import (
	"testing"

	"github.com/ypetya/fitmanager/cmd/mocks"
)

func TestCommandListRemote(t *testing.T) {
	// GIVEN
	ds := mocks.FakeDataStore{LoadReturn: true}
	// WHEN
	commandListRemotes([]string{"dir"}, &ds)
	// THEN
	if ds.Calls[0] != "Load" {
		t.Error("Should load database first")
	}
	if ds.Calls[1] != "ListRemotes" {
		t.Error("Should call ListRemotes second")
	}
}

func TestCommandAddRemote(t *testing.T) {
	// GIVEN
	ds := mocks.FakeDataStore{LoadReturn: true}
	// WHEN
	commandAddRemote([]string{"dir", "new", "/media/asdfadsf"}, &ds)
	// THEN
	if ds.Calls[0] != "Load" {
		t.Error("Should load database first")
	}
	if ds.Calls[1] != "AddRemote" || ds.Passed[0] != "new" {
		t.Error("Should call AddRemote with name 'new'")
	}
	if ds.Calls[2] != "Save" {
		t.Error("Should save data-base")
	}
}

func TestCommandDelRemote(t *testing.T) {
	// GIVEN
	ds := mocks.FakeDataStore{LoadReturn: true}
	// WHEN
	commandDelRemote([]string{"dir", "new"}, &ds)
	// THEN
	if ds.Calls[0] != "Load" {
		t.Error("Should load database first")
	}
	if ds.Calls[1] != "DelRemote" || ds.Passed[0] != "new" {
		t.Error("Should call DelRemote with name 'new'")
	}
	if ds.Calls[2] != "Save" {
		t.Error("Should save data-base")
	}
}
