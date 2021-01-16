package main

import (
	m "github.com/ypetya/fitmanager/models"
)

// Implements IFilter
// returns exercises where HR is missing
// and have exactly one overlap with HR band available!
type FilterHrEnhanceable struct {
}

func (FilterHrEnhanceable) Filter(exercises []m.Exercise) []m.Exercise {
	idMap := map[string]*m.Exercise{}
	candidates := []*m.Exercise{}

	// collect candidates
	for ix, e := range exercises {
		if ref, err := e.GetLocalRef(); err == nil {
			ptr := &exercises[ix]
			idMap[ref] = ptr

			// FIXME make it work over multiple (> 0) -> enhancer should know how to pick second
			if len(ptr.OverlappingIds) == 1 && !ptr.Meta.HasHr() {
				candidates = append(candidates, ptr)
			}
		}
	}
	// checking candidates in with the help of idMap
	good := []m.Exercise{}

nextCandidate:
	for ix, c := range candidates {
		havingHr := 0

		for _, id := range c.OverlappingIds {
			if overlapped, ok := idMap[id]; ok {
				if overlapped.Meta.HasHr() {
					havingHr += 1
					if havingHr > 1 {
						continue nextCandidate
					}
				}
			} else {
				panic("DB is corrupt, overlapping id not found:" + id)
			}
		}
		// won
		if havingHr == 1 {
			good = append(good, *candidates[ix])
		}
	}

	return good
}
