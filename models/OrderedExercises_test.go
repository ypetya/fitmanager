package models

import "testing"

func Test_createOrderedExercises(t *testing.T) {
	// GIVEN
	exercises := []Exercise{
		Exercise{
			Meta: ActivityMetaInfo{
				Activity: "Cycling",
				Start:    12,
				End:      18,
			},
		},
		Exercise{
			Meta: ActivityMetaInfo{
				Activity: "Cycling",
				Start:    11,
				End:      16,
			},
		},
	}
	// WHEN
	ord := createOrderedExercises(&exercises)
	// THEN
	if ord.Len() != 4 {
		t.Error("Expected to have Len == 4!")
	}
	if ord.stack[0].ex != &(exercises[1]) {
		t.Error("There is a problem with the ordering!")
	}
}
func Test_createOrderedExercises_Should_Drop_Invalid_Exercises(t *testing.T) {
	// GIVEN
	exercises := []Exercise{
		Exercise{
			Meta: ActivityMetaInfo{
				Activity: "Cycling",
				Start:    11,
				End:      11,
			},
		},
		Exercise{
			Meta: ActivityMetaInfo{
				Activity: "Cycling",
				Start:    12,
				End:      10,
			},
		},
	}
	// WHEN
	ord := createOrderedExercises(&exercises)
	// THEN
	if ord.Len() != 0 {
		t.Error("Expected to have Len == 0!")
	}
}
func Test_createOrderedExercises_Should_clear_existing_overlaps(t *testing.T) {
	// GIVEN
	exercises := []Exercise{
		Exercise{
			Meta: ActivityMetaInfo{
				Activity: "Cycling",
				Start:    11,
				End:      11,
			},
			OverlappingIds: []string{"modified"},
		},
		Exercise{
			Meta: ActivityMetaInfo{
				Activity: "Cycling",
				Start:    12,
				End:      10,
			},
		},
	}
	// WHEN
	createOrderedExercises(&exercises)
	// THEN
	if len(exercises[0].OverlappingIds) != 0 {
		t.Error("Expected to not have overlaps!")
	}
}
func Test_createOrderedExercises_Should_work_on_nil_exercises(t *testing.T) {
	// GIVEN
	var exercises *[]Exercise
	exercises = nil
	// WHEN
	ord := createOrderedExercises(exercises)
	// THEN
	if ord.Len() != 0 {
		t.Error("Expected to have Len == 0!")
	}
}
