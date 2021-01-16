package main

import (
	"testing"

	"github.com/ypetya/fitmanager/cmd/mocks"
)

func TestCommandList_ShouldListExercises(t *testing.T) {
	// GIVEN
	ds := mocks.FakeDataStore{LoadReturn: true}
	// WHEN
	commandList([]string{"dir"}, &ds)
	// THEN
	if ds.Calls[0] != "Load" {
		t.Error("Should load database first")
	}
	if ds.Calls[1] != "GetExercises" {
		t.Error("Should call GetExercises to list")
	}
}

func TestCommandList_Should_call_filter_by_remote(t *testing.T) {
	// GIVEN
	ds := mocks.FakeDataStore{LoadReturn: true}
	// WHEN
	commandList([]string{"dir", "garmin"}, &ds)
	// THEN
	if ds.Calls[0] != "Load" {
		t.Error("Should load database first")
	}
	if ds.Calls[1] != "Filter" {
		t.Error("Should call Filter to list")
	}
}

func Test_formatCreated_Should_return_formatted_day_time(t *testing.T) {
	// GIVEN
	ts := int64(1607797993)
	expected := "2020-12-12 19:33"
	// WHEN
	res := formatCreated(ts)
	// THEN
	if res != expected {
		t.Errorf("Expected to get '%s', but found %s", expected, res)
	}
}

func Test_formatDuration_Should_return_duration_truncated(t *testing.T) {
	// GIVEN
	uEnd := []int64{0, 1, 60 + 40, 12*60*60 + 60 + 59}
	expected := []string{"0", "1s", "1m40s", "12h1m59s"}
	// WHEN - THEN
	i := 0
	for i < len(uEnd) {
		res := formatDuration(0, uEnd[i])
		if res != expected[i] {
			t.Errorf("Expected to get '%s' for %d, found '%s'", expected[i], uEnd[i], res)
		}
		i++
	}
}
