package models

import "testing"

func Test_DataStore_isExerciseValidForUpload_should_return_true_when_no_overlaps_and_no_such_ref(t *testing.T) {
	// GIVEN
	ex := Exercise{
		Meta: ActivityMetaInfo{
			Samples: 26,
		},
		StoredOn: []AnyRemote{
			AnyRemote{
				Target: LocalDB,
				File: &File{
					Path: "testdata/imagine.FIT",
				},
			},
		},
	}
	// WHEN
	ret := ds.isExerciseValidForUpload(&ex, "garmin")
	// THEN
	if ret != true {
		t.Error("Expected to return true when no overlaps and no ref found")
	}
}
func Test_DataStore_isExerciseValidForUpload_should_return_false_when_one_of_the_overlaps_has_ref(t *testing.T) {
	// GIVEN
	ex1 := Exercise{
		Meta: ActivityMetaInfo{
			Bands: []string{"hr"},
		},
		OverlappingIds: []string{"2.FIT"},
		StoredOn: []AnyRemote{
			AnyRemote{
				Target: LocalDB,
				File: &File{
					Path: "1.FIT",
				},
			},
		},
	}
	ex2 := Exercise{
		Meta: ActivityMetaInfo{
			Samples: 26,
			Bands:   []string{"hr", "cad"},
		},
		OverlappingIds: []string{"1.FIT"},
		StoredOn: []AnyRemote{
			AnyRemote{
				Target: LocalDB,
				File: &File{
					Path: "2.FIT",
				},
			},
			AnyRemote{
				Name:   "garmin",
				Target: GarminConnect,
				Id:     "123",
			},
		},
	}
	ds.initStruct()
	ds.Exercises = []Exercise{ex1, ex2}
	ds.buildRefMap()
	// WHEN
	isValid := ds.isExerciseValidForUpload(&ex1, "garmin")
	// THEN
	if isValid {
		t.Error("Expected to be invalid as the overlap has the remote!")
	}
	// WHEN
	isValid = ds.isExerciseValidForUpload(&ex2, "garmin")
	// THEN
	if isValid {
		t.Error("Expected to be invalid as has the remote!")
	}
}
func Test_DataStore_isExerciseValidForUpload_should_return_true_only_when_this_is_the_selected_from_overlapping_exercises(t *testing.T) {
	// GIVEN
	ex1 := Exercise{
		Meta: ActivityMetaInfo{
			Bands: []string{"hr"},
		},
		OverlappingIds: []string{"2.FIT"},
		StoredOn: []AnyRemote{
			AnyRemote{
				Target: LocalDB,
				File: &File{
					Path: "1.FIT",
				},
			},
		},
	}
	ex2 := Exercise{
		Meta: ActivityMetaInfo{
			Samples: 26,
			Bands:   []string{"hr", "cad"},
		},
		OverlappingIds: []string{"1.FIT"},
		StoredOn: []AnyRemote{
			AnyRemote{
				Target: LocalDB,
				File: &File{
					Path: "2.FIT",
				},
			},
		},
	}
	ds.initStruct()
	ds.Exercises = []Exercise{ex1, ex2}
	ds.buildRefMap()
	// WHEN
	isValid := ds.isExerciseValidForUpload(&ex2, "garmin")
	// THEN
	if !isValid {
		t.Error("Expected to have valid to be true when selected!")
	}
	// WHEN
	isValid = ds.isExerciseValidForUpload(&ex1, "garmin")
	// THEN
	if isValid {
		t.Error("Expected to have valid to be false when it is not the selected!")
	}
}
func Test_DataStore_exerciseOverlapsHasRef_should_return_true_when_one_of_the_overlaps_has_remote_as_ref(t *testing.T) {
	// GIVEN
	ex1 := Exercise{
		Meta: ActivityMetaInfo{
			Bands: []string{"hr"},
		},
		OverlappingIds: []string{"2.FIT"},
		StoredOn: []AnyRemote{
			AnyRemote{
				Target: LocalDB,
				File: &File{
					Path: "1.FIT",
				},
			},
		},
	}
	ex2 := Exercise{
		Meta: ActivityMetaInfo{
			Samples: 26,
			Bands:   []string{"hr", "cad"},
		},
		OverlappingIds: []string{"1.FIT"},
		StoredOn: []AnyRemote{
			AnyRemote{
				Target: LocalDB,
				File: &File{
					Path: "2.FIT",
				},
			},
			AnyRemote{
				Name:   "garmin",
				Target: GarminConnect,
				Id:     "123",
			},
		},
	}
	ds.initStruct()
	ds.Exercises = []Exercise{ex1, ex2}
	ds.buildRefMap()
	// WHEN
	ret := ds.exerciseOverlapsHasRef(&ex1, "garmin")
	// THEN
	if !ret {
		t.Error("Expected to have exerciseOvelapsHasRef to be true when overlapping exercise has the remoteRef!")
	}
}
func Test_DataStore_exerciseOverlapsHasRef_should_return_false_when_none_of_the_overlaps_has_remote_as_ref(t *testing.T) {
	// GIVEN
	ex1 := Exercise{
		Meta: ActivityMetaInfo{
			Bands: []string{"hr"},
		},
		OverlappingIds: []string{"2.FIT"},
		StoredOn: []AnyRemote{
			AnyRemote{
				Target: LocalDB,
				File: &File{
					Path: "1.FIT",
				},
			},
		},
	}
	ex2 := Exercise{
		Meta: ActivityMetaInfo{
			Samples: 26,
			Bands:   []string{"hr", "cad"},
		},
		OverlappingIds: []string{"1.FIT"},
		StoredOn: []AnyRemote{
			AnyRemote{
				Target: LocalDB,
				File: &File{
					Path: "2.FIT",
				},
			},
		},
	}
	ds.initStruct()
	ds.Exercises = []Exercise{ex1, ex2}
	ds.buildRefMap()
	// WHEN
	ret := ds.exerciseOverlapsHasRef(&ex1, "garmin")
	// THEN
	if ret {
		t.Error("Expected to have exerciseOvelapsHasRef to be false!")
	}
}
func Test_DataStore_selectExerciseForExportFromOverlaps_should_return_ex_with_most_bands(t *testing.T) {
	// GIVEN
	ex1 := Exercise{
		Meta: ActivityMetaInfo{
			Bands: []string{"hr"},
		},
		OverlappingIds: []string{"2.FIT"},
		StoredOn: []AnyRemote{
			AnyRemote{
				Target: LocalDB,
				File: &File{
					Path: "1.FIT",
				},
			},
		},
	}
	ex2 := Exercise{
		Meta: ActivityMetaInfo{
			Samples: 26,
			Bands:   []string{"hr", "cad"},
		},
		OverlappingIds: []string{"1.FIT"},
		StoredOn: []AnyRemote{
			AnyRemote{
				Target: LocalDB,
				File: &File{
					Path: "2.FIT",
				},
			},
		},
	}
	ds.initStruct()
	ds.Exercises = []Exercise{ex1, ex2}
	ds.buildRefMap()
	// WHEN
	ret := ds.selectExerciseForExportFromOverlaps(&ex1)
	// THEN
	if !ret.equals(ex2) {
		t.Error("Expected to select the exercise with the most bands")
	}
}

func Test_DataStore_selectExerciseForExportFromOverlaps_should_return_consequently_the_same(t *testing.T) {
	// GIVEN
	ex1 := Exercise{
		Meta:           ActivityMetaInfo{},
		OverlappingIds: []string{"2.FIT"},
		StoredOn: []AnyRemote{
			AnyRemote{
				Target: LocalDB,
				File: &File{
					Path: "1.FIT",
				},
			},
		},
	}
	ex2 := Exercise{
		Meta:           ActivityMetaInfo{},
		OverlappingIds: []string{"1.FIT"},
		StoredOn: []AnyRemote{
			AnyRemote{
				Target: LocalDB,
				File: &File{
					Path: "2.FIT",
				},
			},
		},
	}
	ds.RefIndexMax = 0
	ds.initStruct()
	ds.Exercises = []Exercise{ex1, ex2}
	ds.buildRefMap()
	// WHEN
	// we must pass the internal pointers!
	ret1 := ds.selectExerciseForExportFromOverlaps(&ds.Exercises[0])
	ret2 := ds.selectExerciseForExportFromOverlaps(&ds.Exercises[1])

	// THEN
	r1, _ := ret1.GetLocalRef()
	r2, _ := ret2.GetLocalRef()
	if ret1 != ret2 {
		t.Errorf("Expected to select the same exercise, %s != %s", r1, r2)
	}
	if ret1.RefIndex != 2 ||
		ds.Exercises[0].RefIndex != 1 ||
		ds.Exercises[1].RefIndex != 2 {
		t.Error("Expected to return the mostRecent exercise!")
	}
}
