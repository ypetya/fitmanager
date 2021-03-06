package garminconnect

import (
	m "github.com/ypetya/fitmanager/metadataExtractor"
)

func (GarminConnectConnector) Extract(file string) (Activity string,
	Device string,
	Start int64,
	End int64,
	Samples int64,
	Bands []string,
	Created int64,
) {
	return m.Extract(file)
}
