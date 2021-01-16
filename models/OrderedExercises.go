package models

import "sort"

type poi struct {
	time  int64
	begin bool
	ex    *Exercise
}

type orderedExercises struct {
	stack []poi
}

// Len is the number of elements in the collection.
func (p orderedExercises) Len() int {
	return len(p.stack)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (p orderedExercises) Less(i, j int) bool {
	return p.stack[i].time < p.stack[j].time
}

// Swap swaps the elements with indexes i and j.
func (p orderedExercises) Swap(i, j int) {
	p.stack[i], p.stack[j] = p.stack[j], p.stack[i]
}

// prepare overlapping detection
func createOrderedExercises(e *[]Exercise) orderedExercises {
	// build stack
	stack := []poi{}

	if e != nil {
		for i, v := range *e {
			exPointer := &((*e)[i])
			// remove existing overlaps
			exPointer.clearOverlaps()
			// create start-/end-marks
			start, end :=
				poi{begin: true, time: v.Meta.Start, ex: exPointer},
				poi{begin: false, time: v.Meta.End, ex: exPointer}

			// count in only valid exercises
			if end.time > start.time {
				stack = append(stack, start, end)
			}
		}
	}
	// order
	ret := orderedExercises{
		stack: stack,
	}

	sort.Sort(ret)

	return ret
}
