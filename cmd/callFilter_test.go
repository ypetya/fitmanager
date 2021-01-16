package main

import (
	"testing"

	"github.com/ypetya/fitmanager/cmd/mocks"
)

func Test_callFilter_Should_translate_args_to_filter(t *testing.T) {
	// GIVEN
	ds := mocks.FakeDataStore{LoadReturn: true}
	// WHEN
	callFilter(&ds, []string{"garmin"})
	// THEN
	if ds.FilterPassed == nil {
		t.Error("Should construct Filter")
	}
}

func Test_callFilter_Should_create_filter_for_a_single_remote(t *testing.T) {
	// GIVEN
	ds := mocks.FakeDataStore{LoadReturn: true}
	// WHEN
	callFilter(&ds, []string{"garmin"})
	// THEN
	f := *ds.FilterPassed
	filter := f.(*mocks.FakeFilter)
	c := len(filter.Calls)
	if c != 1 {
		t.Errorf("Expected to have 1 call on filter, found %d, %s", c, filter.Calls)
	}
	if len(filter.Params) != 1 {
		t.Errorf("Expected to have 1 filter param, found %d", len(filter.Params))
	}
	if len(filter.Params[0]) != 1 {
		t.Errorf("Expected to have 1 string in first filter param, found %d",
			len(filter.Params[0]))
	}
	if filter.Params[0][0] != "garmin" {
		t.Errorf("Expected to have 'garmin' as the first filter criteria, found %s",
			filter.Params[0][0])
	}
}

func Test_callFilter_Should_translate_band_filters_of_filter(t *testing.T) {
	// GIVEN
	ds := mocks.FakeDataStore{LoadReturn: true}
	// WHEN
	callFilter(&ds, []string{"garmin", "+hr", "zwift", "-hr"})
	// THEN
	f := *ds.FilterPassed
	filter := f.(*mocks.FakeFilter)
	c := len(filter.Calls)
	if c != 2 {
		t.Errorf("Expected to have 2 calls on filter, found %d, %s", c, filter.Calls)
	}
}
