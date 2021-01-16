package models

import (
	"testing"

	"../connectors"
	mocks "./mocks"
)

func Test_DataStore_Import_Should_return_error_when_no_suchRemote(t *testing.T) {
	// GIVEN
	ds.Load("testdata/export-mixed")
	// WHEN
	err := ds.Import("Invalid")
	// THEN
	if err != ErrorRemoteNotFound {
		t.Error("Expected to return ErrorRemoteNotFound")
	}
}

// it should work with multiple importers
func Test_DS_Import_SupportsMultipleRemotes(t *testing.T) {
	//GIVEN
	importer1 := mocks.CreateFakeConnector([]string{"1", "2", "3.fit", "4.fit"})
	importer2 := mocks.CreateFakeConnector([]string{"e3.FIT", "e4.FIT", "e5.FIT", "e6.FIT"})
	mocks.ResetFakeMetadataExtractorCalls()
	fmde.AlwaysNew = false
	// Metadata with prop "Created"=1, 2 existing in testdata !!!
	fmde.FileCreated = map[string]int64{
		"testdata/remotes-saved/1": 1, "testdata/remotes-saved/2": 2, // <- existing exercises
		"testdata/remotes-saved/3.fit": 3, "testdata/remotes-saved/4.fit": 4, // <- from source1, 2 new
		"testdata/remotes-saved/e3.FIT": 3, "testdata/remotes-saved/e4.FIT": 4,
		"testdata/remotes-saved/e5.FIT": 5, "testdata/remotes-saved/e6.FIT": 6, // <- from source2, 2 new
	}

	ds := DataStore{
		connectors: map[RemoteType]connectors.IConnector{
			GarminConnect: importer1,
			Directory:     importer2,
		},
		metadataExtractor: fmde,
		io:                fio,
	}
	if !ds.Load("testdata/remotes-saved") {
		t.Errorf("Expected to return Load succeeded")
	}
	if len(mocks.FakeCalls) != 0 {
		t.Errorf("There should be 0 calls of the importer, found %d", len(mocks.FakeCalls))
	}
	if len(ds.Exercises) != 2 {
		t.Errorf("Initially there should be 2 exercises, but found %d", len(ds.Exercises))
	}
	// lastSync1, lastSync2 := ds.Remotes["garmin"].LastSync, ds.Remotes["fenix"].LastSync
	// WHEN
	ds.Import("garmin")
	if len(mocks.FakeCalls) != 3 {
		t.Errorf("There should be 3 calls(FetchDiff,Import,Import) of the importer, found %d.", len(mocks.FakeCalls))
	}
	ds.Import("fenix")
	// THEN
	if len(mocks.FakeCalls) != 9 {
		t.Errorf("There should be 9 calls of the importer, found %d", len(mocks.FakeCalls))
		t.Error(mocks.FakeCalls)
	}
	if len(ds.Exercises) != 6 {
		t.Errorf("Finally there should be 6 exercises in the store, found %d", len(ds.Exercises))
	}
}

func Test_DS_Import_FromDir_ShoudSetSourceOnDirectoryImporter(t *testing.T) {
	// GIVEN
	importer := mocks.CreateFakeConnector([]string{"1.fit"})

	ds := DataStore{
		connectors: map[RemoteType]connectors.IConnector{
			Directory: importer},
		metadataExtractor: fmde,
		io:                fio,
	}
	if !ds.Load("testdata/remotes-saved") {
		t.Errorf("Expected to return Load succeeded")
	}
	// WHEN
	ds.Import("fenix")
	// THEN
	if len(mocks.FakeCalls) != 3 {
		t.Error("There should be 3 fake importer calls!")
		t.Error(mocks.FakeCalls)
	}

}

func Test_DS_Import_FromNonExistingRemoteShouldNotCallImporter(t *testing.T) {
	// GIVEN
	importer := mocks.CreateFakeConnector([]string{"1.fit"})

	ds := DataStore{
		connectors: map[RemoteType]connectors.IConnector{
			Directory: importer},
		metadataExtractor: fmde,
		io:                fio,
	}
	if !ds.Load("testdata/remotes-saved") {
		t.Errorf("Expected to return Load succeeded")
	}
	// WHEN
	ds.Import("fenix3")
	// THEN
	if len(mocks.FakeCalls) != 0 {
		t.Error("There should be exactly 0 fake importer calls!")
	}
}

func Test_DS_Import_FromRemote(t *testing.T) {
	// GIVEN
	importer := mocks.CreateFakeConnector([]string{"1", "2", "3", "4"})
	mocks.ResetFakeMetadataExtractorCalls()
	fmde.AlwaysNew = true

	ds := DataStore{
		connectors: map[RemoteType]connectors.IConnector{
			GarminConnect: importer},
		metadataExtractor: fmde,
		io:                fio,
	}
	if !ds.Load("testdata/remotes-saved") {
		t.Errorf("Expected to return Load succeeded")
	}
	if len(ds.Exercises) != 2 {
		t.Errorf("Initially there should be 2 exercises, but found %d", len(ds.Exercises))
	}
	lastSync := ds.Remotes["garmin"].LastSync
	// WHEN
	err := ds.Import("garmin")
	// THEN
	if err != nil {
		t.Error("Should not return error on import call!")
	}
	if len(mocks.FakeCalls) != 3 {
		t.Errorf("There should be 3 calls of the importer")
	}
	if mocks.FakeCalls[0] != "FetchDiff" {
		t.Error("The first call should be FetchDiff")
	}
	if mocks.FakeCalls[1] != "Import:3" {
		t.Error("Should call for the first file")
	}
	if mocks.FakeCalls[2] != "Import:4" {
		t.Error("Should call for the second file")
	}
	if len(mocks.FakeMetadataExtractorCalls) != 2 {
		t.Errorf("Expected to extract metadata from imports; calls: %s",
			mocks.FakeMetadataExtractorCalls)
	}
	if len(ds.Exercises) != 4 {
		t.Errorf("Finally there should be 4 exercises in the store, found %d", len(ds.Exercises))
	}
	newStoredRefs := ds.collectStoredRefsForRemote("garmin")
	if newStoredRefs[2] != "3" {
		t.Error("ExternalId should have been saved for newly imported exercises!")
	}
	if lastSync == ds.Remotes["garmin"].LastSync {
		t.Error("Last sync time-stamp should have been updated for remote!")
	}
	if !fio.SameSize("testdata/remotes-saved/ds-saved-try.json") {
		t.Errorf("The saved metainfo size should match when files are the same!")
	}
}

// Test match checks stored remote!
func Test_DS_Import_FromRemote_ShouldExtendExistingActivity(t *testing.T) {
	// GIVEN
	importer := mocks.CreateFakeConnector([]string{"1", "2", "3", "4"})
	mocks.ResetFakeMetadataExtractorCalls()
	// Metadata with prop "Created"=1, 2 existing in testdata !!!
	fmde.FileCreated = map[string]int64{
		"1": 1, "2": 2, // <- these are existing
		"3": 3, "4": 4, // new exercises!!
	}
	fio := &mocks.FakeIO{}

	ds := DataStore{
		connectors: map[RemoteType]connectors.IConnector{
			GarminConnect: importer},
		metadataExtractor: fmde,
		io:                fio,
	}
	if !ds.Load("testdata/remotes-extend-existing") {
		t.Errorf("Expected to return Load succeeded")
	}
	if len(ds.Exercises) != 2 {
		t.Errorf("Initially there should be 2 exercises, but found %d", len(ds.Exercises))
	}
	// WHEN
	ds.Import("garmin")
	// THEN
	if len(ds.Exercises) != 4 {
		t.Errorf("Finally there should be 4 exercises in the store, found %d", len(ds.Exercises))
	}
	newStoredRefs := ds.collectStoredRefsForRemote("garmin")
	if newStoredRefs[0] != "1" {
		t.Error("ExternalId should have been saved for existing exercise!")
	}
	lastSync := ds.Remotes["garmin"].LastSync
	if ds.Exercises[0].StoredOn[1].LastSync != lastSync {
		t.Errorf("Last sync should have been updated for exercises 0!")
	}
}

func Test_DS_ImportWithoutAnImporter_ShouldReturnError(t *testing.T) {
	// GIVEN
	fmde := mocks.CreateFakeMetadataExtractor()
	ds := DataStore{
		connectors:        map[RemoteType]connectors.IConnector{},
		metadataExtractor: fmde,
	}
	ds.Load("testdata/remotes-saved")
	// WHEN

	if ds.Import("garmin") != ErrorImporterNotFound {
		// THEN
		t.Errorf("Expected to return Error because there is no importer for the remote types!")
	}
}
