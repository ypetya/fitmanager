package models

import "testing"

func Test_DataStore_buildRefMap_should_set_refIndices(t *testing.T) {
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
	// WHEN
	ds.buildRefMap()
	// THEN
	if ds.Exercises[0].RefIndex != 1 {
		t.Errorf("Expected to set refIndex - 1, found %d", ex1.RefIndex)
	}
	if ds.Exercises[1].RefIndex != 2 {
		t.Errorf("Expected to set refIndex - 2, found %d", ex2.RefIndex)
	}
	if ds.RefIndexMax != 2 {
		t.Errorf("Expected to set db.refIndexMax, found %d", ds.RefIndexMax)
	}
}
