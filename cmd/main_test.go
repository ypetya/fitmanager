package main

import (
	"testing"

	"github.com/ypetya/fitmanager/cmd/mocks"
)

func TestDefineCommands(t *testing.T) {
	// GIVEN
	e := mocks.NewFakeCommands()
	// WHEN
	defineCommands(&e)
	// THEN
	if len(e.Calls) != 8 {
		t.Errorf("There should be 8 commands, found %d", len(e.Calls))
	}
	if !e.HasCommand("init") {
		t.Error("Should define command for init")
	}
	if !e.HasCommand("import") {
		t.Error("Should define command for import")
	}
	if !e.HasCommand("add") {
		t.Error("Should define command for addRemote")
	}
	if !e.HasCommand("remove") {
		t.Error("Should define command for removeRemote")
	}
	if !e.HasCommand("remotes") {
		t.Error("Should define command for listRemote")
	}
	if !e.HasCommand("summary") {
		t.Error("Should define command for summary")
	}
	if !e.HasCommand("list") {
		t.Error("Should define command for listing exercises")
	}
	if !e.HasCommand("export") {
		t.Error("Should define command for export")
	}
}
