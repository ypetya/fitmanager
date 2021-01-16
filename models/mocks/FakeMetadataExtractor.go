package mocks

import (
	m "github.com/ypetya/fitmanager/metadataExtractor"
)

var FakeMetadataExtractorCalls []string

type FakeMetadataExtractor struct {
	AlwaysNew     bool
	MockProviders map[string]m.IMetadataExtractor
	FileCreated   map[string]int64
}

func (f FakeMetadataExtractor) Extract(fn string) (string, string, int64, int64, int64, []string, int64) {
	FakeMetadataExtractorCalls = append(FakeMetadataExtractorCalls, fn)
	if f.AlwaysNew {
		return "Cycling", "", 0, 0, 0, []string{}, int64(len(FakeMetadataExtractorCalls))
	}
	if cb, ok := f.MockProviders[fn]; ok {
		return cb.Extract(fn)
	}
	if v, ok := f.FileCreated[fn]; ok {
		return "Cycling", "", 0, 0, 0, []string{}, v
	}
	return "Cycling", "", 0, 0, 0, []string{}, 0
}

func CreateFakeMetadataExtractor() FakeMetadataExtractor {
	ResetFakeMetadataExtractorCalls()
	fmde := FakeMetadataExtractor{}
	fmde.MockProviders = make(map[string]m.IMetadataExtractor)
	fmde.FileCreated = make(map[string]int64)
	return fmde
}

func ResetFakeMetadataExtractorCalls() {
	FakeMetadataExtractorCalls = []string{}
}
