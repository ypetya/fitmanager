package models

import (
	"../connectors"
	"../internal"
	"../metadataExtractor"
)

const VERSION = 8

func NewDataStore(
	connectors map[RemoteType]connectors.IConnector,
	extractor metadataExtractor.IMetadataExtractor,
	enhancers []Enhancer,
) *DataStore {
	io := internal.EffectiveIO{}

	db := DataStore{
		Version:           VERSION,
		io:                io,
		metadataExtractor: extractor,
		connectors:        connectors,
		enhancers:         enhancers,
	}
	return &db
}
