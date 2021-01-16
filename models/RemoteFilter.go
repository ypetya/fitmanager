package models

type RemoteFilter struct {
	// the first string in a condition is remote name
	// the next strings are conditions on meta data
	// for example [ "garmin" "+hr" ... ]
	// or negating [ "garmin" "-hr" ]
	remoteConditions [][]string
}

func (f *RemoteFilter) Remote(remote string, condition []string) {
	cond := make([]string, len(condition)+1)
	cond[0] = remote
	copy(cond[1:], condition)
	f.remoteConditions = append(f.remoteConditions, cond)
}

func (f *RemoteFilter) Filter(exercises []Exercise) []Exercise {
	var found []Exercise

	if len(f.remoteConditions) > 0 {
		found = make([]Exercise, 0, cap(exercises))
		for ix, ex := range exercises {
			for _, conditionList := range f.remoteConditions {
				if len(conditionList) == 1 {
					remoteName := conditionList[0]
					if _, err := ex.GetRemoteRef(remoteName); err == nil {
						found = append(found, exercises[ix])
					}
				} else {
					panic("Band filer not implemented!")
				}
			}
		}
	} else {
		found = make([]Exercise, len(exercises))

		copy(found, exercises)
	}

	return found
}
