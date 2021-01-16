package models

// TODO add hasBand function
type ActivityMetaInfo struct {
	// The type of activitiy recorded
	Activity string `json:"activity,omitempty"`
	// Device used for recording the exercise
	Device string `json:"device,omitempty"`
	// begining time when the exercise started
	Start int64 `json:"start,omitempty"`
	// finish time when the exercise ended
	End int64 `json:"end,omitempty"`
	// count of samples stored
	Samples int64 `json:"samples,omitempty"`
	// bands stored in the sample
	// FIXME order - assuming MetaDataExtractor gives the same always, no manual edit!
	Bands []string `json:"bands,omitempty"`
	// time-stamp when the meta-info parsed
	Created int64 `json:"created,omitempty"`
}

func (a ActivityMetaInfo) HasHr() bool {
	for _, v := range a.Bands {
		if v == "hr" {
			return true
		}
	}
	return false
}

func (a ActivityMetaInfo) equals(o IEquals) bool {
	switch o.(type) {
	case ActivityMetaInfo:
		s := o.(ActivityMetaInfo)
		if len(a.Bands) == len(s.Bands) {
			// TODO assuming ordered!
			for i, _ := range a.Bands {
				if a.Bands[i] != s.Bands[i] {
					return false
				}
			}
		} else {
			return false
		}

		return a.Activity == s.Activity &&
			a.Start == s.Start &&
			a.End == s.End &&
			a.Created == s.Created &&
			a.Samples == s.Samples
	default:
		return false
	}
}
