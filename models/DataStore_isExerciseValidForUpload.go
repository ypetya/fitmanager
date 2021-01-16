package models

// TODO simplify
func (db *DataStore) isExerciseValidForUpload(ex *Exercise, remoteName string) bool {
	if _, err := ex.GetRemoteRef(remoteName); err == nil {
		return false
	}
	// when no overlaps and no remote ref
	if len(ex.OverlappingIds) == 0 {
		if _, err := ex.GetRemoteRef(remoteName); err == ErrorNoSuchRef {
			return true
		}
		return false
	}
	// or when overlaps and no remote-ref and is the selected one!
	if !db.exerciseOverlapsHasRef(ex, remoteName) {
		selectedOne := db.selectExerciseForExportFromOverlaps(ex)
		if selectedOne == ex {
			return true
		}
	}

	// otherwise
	return false
}

// returns true when at least one from overlapping exercises has "remote" ref
func (db *DataStore) exerciseOverlapsHasRef(ex *Exercise, remoteName string) bool {
	for _, id := range ex.OverlappingIds {
		if e, ok := db.refMap[id]; ok {
			if _, err := e.GetRemoteRef(remoteName); err == nil {
				return true
			}
		} else {
			panic("RefMap is not up-to-date! missing exercise id:" + id)
		}
	}
	return false
}

func (db *DataStore) selectExerciseForExportFromOverlaps(ex *Exercise) *Exercise {
	if len(ex.OverlappingIds) == 0 {
		return ex
	}
	maxBands := len(ex.Meta.Bands)
	maxRefIx := ex.RefIndex
	var selByBands, selByRefIx *Exercise
	selByBands, selByRefIx = ex, ex

	for _, id := range ex.OverlappingIds {
		e := db.refMap[id]
		newMax := len(e.Meta.Bands)
		if newMax > maxBands {
			maxBands = newMax
			selByBands = e
		}
		if e.RefIndex > maxRefIx {
			maxRefIx = e.RefIndex
			selByRefIx = e
		}
	}
	if len(selByRefIx.Meta.Bands) < maxBands {
		return selByBands
	}
	return selByRefIx
}
