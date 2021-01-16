package garminconnect

import "github.com/ypetya/fitmanager/internal"

func NewGarminConnectConnector() *GarminConnectConnector {
	c := newClient()
	io := internal.EffectiveIO{}
	return &GarminConnectConnector{
		authenticator: c,
		exporter:      c,
		fetcher:       c,
		io:            io,
	}
}
