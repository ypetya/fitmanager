package models

import "testing"

var ex1 Exercise
var ex2 Exercise
var exWithoutRemote Exercise
var exWithRemote Exercise
var exs []Exercise

func init() {

	ex1 = Exercise{
		Meta: ActivityMetaInfo{
			Device:  "GARMIN",
			Samples: 26,
			Bands:   []string{"HR"},
		},
		StoredOn: []AnyRemote{
			AnyRemote{
				Target: LocalDB,
				File: &File{
					Path: "testdata/imagine.FIT",
				},
			},
			AnyRemote{
				Name:   "garmin",
				Target: GarminConnect,
				Id:     "123",
			},
		},
	}
	ex2 = Exercise{
		Meta: ActivityMetaInfo{
			Device:  "GARMIN",
			Samples: 26,
			Bands:   []string{"HR"},
		},
		StoredOn: []AnyRemote{
			AnyRemote{
				Target: LocalDB,
				File: &File{
					Path: "testdata/imagine2.FIT",
				},
			},
		},
	}
	exWithoutRemote = Exercise{
		Meta: ActivityMetaInfo{
			Device:  "GARMIN",
			Samples: 26,
			Bands:   []string{"HR"},
		},
	}

	exs = []Exercise{
		Exercise{
			Meta: ActivityMetaInfo{
				Device:  "GARMIN",
				Samples: 27,
				Bands:   []string{"HR"},
			},
		},
		Exercise{
			Meta: ActivityMetaInfo{
				Device:  "GARMIN",
				Samples: 26,
				Bands:   []string{"HR", "Cadence"},
			},
		},
		Exercise{
			Meta: ActivityMetaInfo{
				Device:  "GARMIN",
				Samples: 26,
				Bands:   []string{"Cadence"},
			},
		},
		Exercise{
			Meta: ActivityMetaInfo{
				Device:  "GARMIN",
				Samples: 26,
			},
		},
	}
}

func Test_addOverlap_Should_AddLocalDBRef(t *testing.T) {
	// GIVEN
	ex1.OverlappingIds = []string{}
	ex2.OverlappingIds = []string{}
	// WHEN
	ex2.addOverlap(&ex1)
	// THEN
	if len(ex2.OverlappingIds) != 1 {
		t.Error("Did not add id to OverlappingIds!")
	}
	if len(ex1.OverlappingIds) != 1 {
		t.Error("Did not add id to OverlappingIds - other exercise!")
	}
}
func Test_addOverlap_Should_Return_Added_Ovelaps_Count(t *testing.T) {
	// GIVEN
	ex1.OverlappingIds = []string{}
	ex2.OverlappingIds = []string{}
	// WHEN
	ret := ex2.addOverlap(&ex1)
	// THEN
	if ret != 2 {
		t.Error("First calll should find both overlapping, returning 2")
	}
	// WHEN
	ret = ex1.addOverlap(&ex2)
	// THEN
	if ret != 0 {
		t.Error("Did not return 0!")
	}
	// WHEN
	ret = ex2.addOverlap(&ex1)
	// THEN
	if ret != 0 {
		t.Error("Did not return 0!")
	}
}

func Test_addOverlapId_Should_AddOverlapOnce(t *testing.T) {
	// GIVEN
	ex1.OverlappingIds = []string{}
	// WHEN
	k := ex1.addOverlapId("1")
	l := ex1.addOverlapId("1")
	// THEN
	if !k {
		t.Errorf("Should return true for the first call!")
	}
	if l {
		t.Errorf("Should return false, overlappingIds can't duplicate!")
	}
}

func Test_Equals_ForSameExercise(t *testing.T) {
	if !ex1.equals(exWithoutRemote) {
		t.Error("Expected to be equals!")
	}
}

func Test_Diff_Exercises(t *testing.T) {
	for i, e := range exs {
		if ex1.equals(e) {
			t.Errorf("Expected to be NOT equals! index: %d", i)
		}
	}
}

func Test_Exercise_GetremoteRef_ShouldReturnRemoteId(t *testing.T) {
	ref, err := ex1.GetRemoteRef("garmin")
	if err != nil {
		t.Error("Expected not to return an error when has a remote store with the speciefied name")
	}
	if ref != "123" {
		t.Errorf("Expected to return the stored AcytivityId, got %s", ref)
	}
}
func Test_Exercise_GetLocalRef_ShouldReturnLocalFile(t *testing.T) {
	ref, err := ex1.GetLocalRef()
	if err != nil {
		t.Error("Expected not to return an error when store has a localDB entry")
	}
	if ref != "testdata/imagine.FIT" {
		t.Errorf("Expected to return the stored AcytivityId, got %s", ref)
	}
}
func Test_Exercise_GetremoteRef_ShouldReturnErrorWhenNoRemoteWithName(t *testing.T) {
	ref, err := ex1.GetRemoteRef("garminX")
	if ref != "" {
		t.Error("Expected to return a ref == nil when remote is missing!")
	}
	if err != ErrorNoSuchRef {
		t.Error("Expected to return an error when remote is missing")
	}
}
