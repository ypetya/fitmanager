package garminconnect

import "../../internal"

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
