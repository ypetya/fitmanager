package models

/**
  Metadata of an excercise stored in the local database
*/
type Exercise struct {
	// parsed meta info
	Meta ActivityMetaInfo `json:"meta,omitempty"`

	// places where the samples are located
	StoredOn []AnyRemote `json:"storedOn,omitempty"`

	// when the exercise start/end time is overlapping another exercise, it's ids are stored
	OverlappingIds []string `json:"overlaps,omitempty"`

	// refIndex is set by buildRefMap
	// being able to provide reproducable ordering
	RefIndex int `json:"x,omitempty"`
}

// Implements IEquals
func (e Exercise) equals(o IEquals) bool {
	switch o.(type) {
	case Exercise:
		return e.Meta.equals(o.(Exercise).Meta)
	default:
		return false
	}
}

func (e *Exercise) GetLocalRef() (string, error) {
	for _, r := range e.StoredOn {
		if r.Target == LocalDB {
			return r.GetRef(), nil
		}
	}
	return "", ErrorNoSuchRef
}

// Implements IRemotesAggregator
func (e *Exercise) GetRemoteRef(remoteName string) (string, error) {
	for _, r := range e.StoredOn {
		if r.Name == remoteName {
			return r.GetRef(), nil
		}
	}

	return "", ErrorNoSuchRef
}

func (e *Exercise) clearOverlaps() {
	e.OverlappingIds = []string{}
}

// Registers overlap for both exercises
// returns the count of the newly registered overlaps
func (e *Exercise) addOverlap(o *Exercise) int {
	count := 0
	if id, err := o.GetLocalRef(); err == nil {
		if e.addOverlapId(id) {
			count++
		}
	}
	if id, err := e.GetLocalRef(); err == nil {
		if o.addOverlapId(id) {
			count++
		}
	}
	return count
}

// Appends a new overrlapId, returns true when successful
// false otherwise (alreadt existing id)
func (e *Exercise) addOverlapId(newId string) bool {
	for _, r := range e.OverlappingIds {
		if r == newId {
			return false
		}
	}
	e.OverlappingIds = append(e.OverlappingIds, newId)

	return true
}
