package models

import "testing"

func Test_Should_Detect_overlapping_exercises_1(t *testing.T) {
	// GIVEN
	ds := DataStore{
		Exercises: []Exercise{
			Exercise{
				Meta: ActivityMetaInfo{
					Activity: "Cycling",
					Start:    12,
					End:      18,
				},
				StoredOn: []AnyRemote{
					AnyRemote{
						Target: LocalDB,
						File:   &File{Path: "1.fit"},
					},
				},
			},
			Exercise{
				Meta: ActivityMetaInfo{
					Activity: "Cycling",
					Start:    11,
					End:      16,
				},
				StoredOn: []AnyRemote{
					AnyRemote{
						Target: LocalDB,
						File:   &File{Path: "2.fit"},
					},
				},
			},
		},
	}
	// WHEN
	ds.markOverlappingExercises()
	// THEN
	if ref, _ := ds.Exercises[0].GetLocalRef(); ref != "1.fit" {
		t.Error("No ref for exercise 0")
	}
	lo := len(ds.Exercises[0].OverlappingIds)
	if lo != 1 {
		t.Errorf("Expected to have overlappingIds length == 1, found %d!", lo)
	}
	ex0_o := ds.Exercises[0].OverlappingIds[0]
	if ex0_o != "2.fit" {
		t.Errorf("Expected to have an overlap with id == '2.fit', found %s!", ex0_o)
	}
	if len(ds.Exercises[1].OverlappingIds) != 1 {
		t.Error("Expected to have overlappingIds length == 1!")
	}
	if ds.Exercises[1].OverlappingIds[0] != "1.fit" {
		t.Error("Expected to have an overlap with id == '1.fit'!")
	}
}
func Test_Should_Detect_overlapping_exercises_2(t *testing.T) {
	// GIVEN
	ds := DataStore{
		Exercises: []Exercise{
			Exercise{
				Meta: ActivityMetaInfo{
					Activity: "Cycling",
					Start:    11,
					End:      15,
				},
				StoredOn: []AnyRemote{
					AnyRemote{
						Target: LocalDB,
						File:   &File{Path: "1.fit"},
					},
				},
			},
			Exercise{
				Meta: ActivityMetaInfo{
					Activity: "Cycling",
					Start:    13,
					End:      16,
				},
				StoredOn: []AnyRemote{
					AnyRemote{
						Target: LocalDB,
						File:   &File{Path: "2.fit"},
					},
				},
			},
		},
	}
	// WHEN
	ds.markOverlappingExercises()
	// THEN
	if ref, _ := ds.Exercises[0].GetLocalRef(); ref != "1.fit" {
		t.Error("No ref for exercise 0")
	}
	lo := len(ds.Exercises[0].OverlappingIds)
	if lo != 1 {
		t.Errorf("Expected to have overlappingIds length == 1, found %d!", lo)
	}
	ex0_o := ds.Exercises[0].OverlappingIds[0]
	if ex0_o != "2.fit" {
		t.Errorf("Expected to have an overlap with id == '2.fit', found %s!", ex0_o)
	}
	if len(ds.Exercises[1].OverlappingIds) != 1 {
		t.Error("Expected to have overlappingIds length == 1!")
	}
	if ds.Exercises[1].OverlappingIds[0] != "1.fit" {
		t.Error("Expected to have an overlap with id == '1.fit'!")
	}
}
func Test_markOvelappingExercises_Should_return_count(t *testing.T) {
	// GIVEN
	ds := DataStore{
		Exercises: []Exercise{
			Exercise{
				Meta: ActivityMetaInfo{
					Activity: "Cycling",
					Start:    12,
					End:      10,
				},
				StoredOn: []AnyRemote{
					AnyRemote{
						Target: LocalDB,
						File:   &File{Path: "1.fit"},
					},
				},
			},
			Exercise{
				Meta: ActivityMetaInfo{
					Activity: "Cycling",
					Start:    13,
					End:      15,
				},
				StoredOn: []AnyRemote{
					AnyRemote{
						Target: LocalDB,
						File:   &File{Path: "2.fit"},
					},
				},
			},
			Exercise{
				Meta: ActivityMetaInfo{
					Activity: "Cycling",
					Start:    14,
					End:      16,
				},
				StoredOn: []AnyRemote{
					AnyRemote{
						Name:   "local",
						Target: LocalDB,
						File:   &File{Path: "3.fit"},
					},
				},
			},
		},
	}
	// WHEN
	ret := ds.markOverlappingExercises()
	// THEN
	if ret != 2 {
		t.Errorf("Expected to have 2 overlappingIds, found %d!", ret)
	}
}
