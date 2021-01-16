package main

import (
	"testing"

	m "github.com/ypetya/fitmanager/models"
)

var filterTestData []m.Exercise

func init() {
	filterTestData = []m.Exercise{
		m.Exercise{
			Meta: m.ActivityMetaInfo{
				Device:  "GARMIN",
				Samples: 27,
				Bands:   []string{"cad", "pos"},
			},
			OverlappingIds: []string{"1.FIT"},
			StoredOn: []m.AnyRemote{
				m.AnyRemote{
					Target: m.LocalDB,
					File: &m.File{
						Path: "imagine.FIT",
					},
				},
			},
		},
		m.Exercise{
			Meta: m.ActivityMetaInfo{
				Device:  "GARMIN",
				Samples: 26,
				Bands:   []string{"hr", "cad"},
			},
			OverlappingIds: []string{"imagine.fit"},
			StoredOn: []m.AnyRemote{
				m.AnyRemote{
					Target: m.LocalDB,
					File: &m.File{
						Path: "1.FIT",
					},
				},
			},
		},
		m.Exercise{
			Meta: m.ActivityMetaInfo{
				Device:  "GARMIN",
				Samples: 26,
				Bands:   []string{"cad"},
			},
			OverlappingIds: []string{"12.fit", "13.fit"},
			StoredOn: []m.AnyRemote{
				m.AnyRemote{
					Target: m.LocalDB,
					File: &m.File{
						Path: "11.fit",
					},
				},
			},
		},
		m.Exercise{
			Meta: m.ActivityMetaInfo{
				Device:  "GARMIN",
				Samples: 26,
				Bands:   []string{"cad"},
			},
			OverlappingIds: []string{"11.fit", "13.fit"},
			StoredOn: []m.AnyRemote{
				m.AnyRemote{
					Target: m.LocalDB,
					File: &m.File{
						Path: "12.fit",
					},
				},
			},
		},
		m.Exercise{
			Meta: m.ActivityMetaInfo{
				Device:  "GARMIN",
				Samples: 26,
				Bands:   []string{"cad"},
			},
			OverlappingIds: []string{"12.fit", "11.fit"},
			StoredOn: []m.AnyRemote{
				m.AnyRemote{
					Target: m.LocalDB,
					File: &m.File{
						Path: "13.fit",
					},
				},
			},
		},
		m.Exercise{
			Meta: m.ActivityMetaInfo{
				Device:  "GARMIN",
				Samples: 26,
			},
		},
	}
}

func Test_filter_Should_return_only_enhanceable_exercises(t *testing.T) {
	// GIVEN
	filter := FilterHrEnhanceable{}
	// WHEN
	ret := filter.Filter(filterTestData)
	// THEN
	if len(ret) != 1 {
		t.Errorf("Should return matching exercise only! expected 1 only,\n found  %d",
			len(ret))
	}
}
