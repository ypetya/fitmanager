package models

func (db *DataStore) markOverlappingExercises() int {
	ordered := createOrderedExercises(&db.Exercises)
	count := 0

	ongoing := []*Exercise{}
	for j, v := range ordered.stack {
		lo := len(ongoing)
		if v.begin {
			if lo > 0 {
				// set foreign ids
				for k, _ := range ongoing {
					// register the new overlaps
					count += ongoing[k].addOverlap(ordered.stack[j].ex)
				}
			}
			// remember started
			ongoing = append(ongoing, ordered.stack[j].ex)
		} else {
			// forget stopped
			k := 0
			// which one has ended?
			for k, _ = range ongoing {
				if ongoing[k] == v.ex {
					break
				}
			}
			// remove from slice
			switch k {
			case 0: // first
				{
					ongoing = ongoing[1:]
				}
			case lo - 1: // last
				{
					ongoing = ongoing[:lo-1]
				}
			default: // in the middle
				{
					ongoing = append(ongoing[:k], ongoing[k+1:]...)
				}
			}
		}
	}
	return count
}
