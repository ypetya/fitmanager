package models

import "testing"

var remoteFilterTestData []Exercise

func init() {
	remoteFilterTestData = []Exercise{
		Exercise{
			Meta: ActivityMetaInfo{
				Device:  "GARMIN",
				Samples: 27,
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
		},
		Exercise{
			Meta: ActivityMetaInfo{
				Device:  "GARMIN",
				Samples: 26,
				Bands:   []string{"HR", "Cadence"},
			},
			StoredOn: []AnyRemote{
				AnyRemote{
					Name:   "garmin",
					Target: GarminConnect,
					Id:     "1.FIT",
				},
			},
		},
		Exercise{
			Meta: ActivityMetaInfo{
				Device:  "GARMIN",
				Samples: 26,
				Bands:   []string{"Cadence"},
			},
			StoredOn: []AnyRemote{
				AnyRemote{
					Name:   "other",
					Target: Directory,
					Id:     "smthg.fit",
				},
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

func Test_filter_Should_return_all_exercises_by_default(t *testing.T) {
	// GIVEN
	filter := RemoteFilter{}
	// WHEN
	ret := filter.Filter(remoteFilterTestData)
	// THEN
	if len(ret) != len(remoteFilterTestData) {
		t.Errorf("Should return all the exercises! len(exs)=%d, got %d",
			len(remoteFilterTestData),
			len(ret))
	}
}

func Test_filter_Should_return_filtered_exercises_for_a_remote(t *testing.T) {
	// GIVEN
	filter := RemoteFilter{
		remoteConditions: [][]string{
			[]string{"garmin"},
		},
	}
	// WHEN
	ret := filter.Filter(remoteFilterTestData)
	// THEN
	if len(ret) != 1 {
		t.Errorf("Should return 1 exercise! got %d", len(ret))
	}

	if !ret[0].equals(remoteFilterTestData[1]) {
		t.Error("The found exercise did not match!")
	}
}
