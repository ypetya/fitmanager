package main

import (
	"testing"

	"./mocks"
)

func TestCommandSummary_should_get_exercises(t *testing.T) {
	// GIVEN
	ds := mocks.FakeDataStore{LoadReturn: true}
	// WHEN
	commandSummary([]string{"dir"}, &ds)
	// THEN
	if ds.Calls[0] != "Load" {
		t.Error("Should load database first")
	}
	if ds.Calls[1] != "GetExercises" {
		t.Error("Should call GetExercises to apply calculations on")
	}
}
