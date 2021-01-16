package models

import (
	"strconv"
	"testing"

	"../connectors"
	"./mocks"
)

type FakeFilter struct{}

func (FakeFilter) Filter(exercises []Exercise) []Exercise {
	mocks.LogCall("FakeFilter:" + strconv.Itoa(len(exercises)))
	ret := []Exercise{}
	for _, e := range exercises {
		if e.Meta.Samples == 27 {
			ret = append(ret, e)
			break
		}
	}
	return ret
}

var ds DataStore

func beforeAll_Export_Suite() {
	mocks.ResetFakeMetadataExtractorCalls()

	conn := mocks.CreateFakeConnector([]string{"1", "2", "3.fit", "4.fit"})

	ds = DataStore{
		io:                fio,
		connectors:        map[RemoteType]connectors.IConnector{Directory: conn},
		metadataExtractor: fmde,
		enhancers: []Enhancer{
			Enhancer{
				Filter:   FakeFilter{},
				Function: mocks.FakeEnhancer{},
				Name:     "hr",
			},
		},
	}
}

func Test_DataStore_addExerciseRemote_should_add_remote_to_exercise(t *testing.T) {
	// GIVEN
	ds.Load("testdata/export-mixed")
	remoteTemplate := ds.Remotes["ex"]
	remote := remoteTemplate.Fill("newId")
	ep := &ds.Exercises[0]
	// WHEN
	ds.addExerciseRemote(ep, remote)
	// THEN
	if len(ds.Exercises[0].StoredOn) < 3 {
		t.Error("Expected to add remote!")
	}
}

func Test_DataStore_Export_Should_return_error_when_no_suchRemote(t *testing.T) {
	// GIVEN
	ds.Load("testdata/export-mixed")
	// WHEN
	err := ds.Export("Invalid")
	// THEN
	if err != ErrorRemoteNotFound {
		t.Error("Expected to return ErrorRemoteNotFound")
	}
}

type MetaDataExtractorForNewFile struct{}

func (MetaDataExtractorForNewFile) Extract(file string) (Activity string,
	// Device used for recording the exercise
	Device string,
	// begining time when the exercise started
	Start int64,
	// finish time when the exercise ended
	End int64,
	// count of samples stored
	Samples int64,
	// bands stored in the sample
	// TODO order
	Bands []string,
	// time-stamp when the meta-info parsed
	Created int64) {
	return "Cycling", "GARMIN", 1, 16, 28, []string{"hr", "cad", "pos"}, 16
}

func Test_DataStore_Export_Should_run_enhancers(t *testing.T) {
	// GIVEN
	mocks.FakeCalls = []string{}

	fmde.AlwaysNew = false
	fmde.MockProviders["testdata/export-mixed/imagine.FIT.hr.fit"] = MetaDataExtractorForNewFile{}

	ds.metadataExtractor = fmde
	// ds - with 4 distinct time exercises, overlapping
	// #1 -> 2 exercises stored - (one HR) -> ready for enhancement
	// #2 -> 3 exercises stored - no HR!
	// #3 -> a single exercise without remote -> testing corner cases
	// #4 -> already exported to remote "ex"

	// WHEN
	ds.Load("testdata/export-mixed")
	succ := ds.Export("ex")
	// THEN
	if succ != nil {
		t.Error("Expected to return nil- meaning, succesful export, found %w", succ)
	}
	// we expect calls:
	// setTarget,
	// Filter,
	// Enhance
	// fetchDiff
	// + 3 export calls for exercises not present on remote
	// ( exporting 1 per overlapping )
	if len(mocks.FakeCalls) != 1+1+1+1+3 {
		t.Errorf("Expected to have 7 calls, found %d, %s",
			len(mocks.FakeCalls),
			mocks.FakeCalls)
	}
	exportCalls := 0

	for _, v := range mocks.FakeCalls {
		if v[0:6] == "Export" {
			exportCalls += 1
		}
	}
	if exportCalls != 3 {
		t.Errorf("Expected to have 3 export calls, found %d", exportCalls)
	}
}
func Test_DataStore_Export_Should_run_on_empty_store(t *testing.T) {
	// GIVEN
	mocks.FakeCalls = []string{}
	exporter := mocks.CreateFakeConnector([]string{"1", "2", "3.fit", "4.fit"})
	connectors := map[RemoteType]connectors.IConnector{Directory: exporter}
	ds := DataStore{
		io: fio,
		Remotes: map[string]AnyRemote{
			"ex": AnyRemote{
				Name:     "someexternalfolder",
				LastSync: 0,
				Target:   Directory,
				File:     &File{Path: "ex"},
			},
		},
		connectors:        connectors,
		metadataExtractor: fmde,
		enhancers: []Enhancer{
			Enhancer{
				Filter:   FakeFilter{},
				Function: mocks.FakeEnhancer{},
				Name:     "hr",
			},
		},
		initialized: true,
	}
	// WHEN
	succ := ds.Export("ex")
	// THEN
	if succ != nil {
		t.Errorf("Expected to return nil- meaning, succesful export, found %w", succ)
	}
	if len(mocks.FakeCalls) != 0 {
		t.Errorf("Expected to have 0 calls. Found:: %d, calls were: %s",
			len(mocks.FakeCalls),
			mocks.FakeCalls)
	}
}
func Test_DataStore_Export_Should_call_exporter(t *testing.T) {
	// GIVEN
	mocks.FakeCalls = []string{}

	// ds - with 4 distinct time exercises, overlapping
	// #1 -> 2 exercises stored - (one HR) -> ready for enhancement
	// #2 -> 3 exercises stored - no HR!
	// #3 -> a single exercise without remote -> testing corner cases
	// #4 -> already exported to remote "ex"

	// WHEN
	ds.Load("testdata/export-mixed")
	ds.Export("ex")
	// THEN
	if mocks.FakeCalls[2][0:7] != "Enhance" {
		t.Errorf("Expected to call enhance! Calls were %s", mocks.FakeCalls)
	}
}

func Test_DataStore_enhance_Should_create_new_Exercise(t *testing.T) {
	// GIVEN

	// ds - with 4 distinct time exercises, overlapping
	// #1 -> 2 exercises stored - (one HR) -> ready for enhancement
	// #2 -> 3 exercises stored - no HR!
	// #3 -> a single exercise without remote -> testing corner cases
	// #4 -> already exported to remote "ex"

	mocks.FakeCalls = []string{}

	// WHEN
	succ := ds.Load("testdata/export-mixed")
	if !succ {
		t.Errorf("Expected to load exercises")
	}
	err := ds.enhance()
	if err != nil {
		t.Errorf("Expected to enhance exercises, found error %w", err)
	}
	// THEN
	if len(ds.Exercises) != 8 {
		t.Errorf("Expect to creat a new exercise!")
	}
	if len(ds.refMap) != 8 {
		t.Errorf("Expected to add ref to the new exercise!")
	}
	if len(mocks.FakeCalls) != 2 {
		t.Errorf("Expect to have 2 mock calls, had %d", len(mocks.FakeCalls))
		t.Error(mocks.FakeCalls)
	}

	if mocks.FakeCalls[0] != "FakeFilter:7" {
		t.Errorf("The first call should be Filter, found %s", mocks.FakeCalls[0])
	}
	if mocks.FakeCalls[1][0:7] != "Enhance" {
		t.Errorf("The second call should be Enhance, found %s", mocks.FakeCalls[1])
	}
}
func Test_DataStore_Export_Should_use_fetchDiff(t *testing.T) {
	// GIVEN
	mocks.FakeCalls = []string{}

	// ds - with 4 distinct time exercises, overlapping
	// #1 -> 2 exercises stored - (one HR) -> ready for enhancement
	// #2 -> 3 exercises stored - no HR!
	// #3 -> a single exercise without remote -> testing corner cases
	// #4 -> already exported to remote "ex"

	// WHEN
	ds.Load("testdata/export-mixed")
	ds.Export("ex")
	// THEN
	if mocks.FakeCalls[6] != "FetchDiff" {
		t.Errorf("Expected to call Fetchdiff! Calls were %s", mocks.FakeCalls)
	}
}

var exporterCalls []string

// removed file exporter
type updater struct{}

func (updater) FetchDiff(k []string) []string {
	return []string{"123.fit"}
}
func (updater) Update(dir string, fileName string, remoteFileName string) error {
	exporterCalls = append(exporterCalls, "Update", dir, fileName, remoteFileName)
	return nil
}

func Test_DataStore_exportRemovedExercises_Should_export_exercises_of_diff(t *testing.T) {
	// GIVEN
	exporterCalls = []string{}
	e := updater{}
	remoteName := "ex"
	ds.Load("testdata/export-mixed")
	// WHEN
	ds.exportRemovedExercises(e, remoteName)
	// THEN
	var found bool
	for _, param := range exporterCalls {
		if param == "Update" {
			found = true
		}
	}
	if !found {
		t.Errorf("Update expected to have been called! %s", exporterCalls)
	}
	if exporterCalls[2] != "123.fit" {
		t.Errorf("Expected to pass fileName to exporter, found %s", exporterCalls[2])
	}
}
func Test_DataStore_updateExercise_Should_remove_obsolete_ref(t *testing.T) {
	// GIVEN
	ds.Load("testdata/export-mixed")
	newEx := Exercise{
		Meta: ds.Exercises[0].Meta,
		StoredOn: []AnyRemote{
			AnyRemote{
				Target: LocalDB,
				File:   &File{Path: "newRef.FIT"}}}}
	if ref, _ := ds.Exercises[0].GetLocalRef(); ref != "imagine.FIT" {
		t.Error("Prerequisite failed")
	}
	// WHEN
	ds.updateExercise(0, &newEx)
	// THEN
	if _, ok := ds.refMap["imagine.FIT"]; ok {
		t.Error("Expected to remove obsolete ref")
	}
}
func Test_DataStore_collectExercisesForRemote_should_build_map(t *testing.T) {
	// GIVEN
	ds.Load("testdata/export-mixed")
	// WHEN
	ids, refMap := ds.collectExercisesForRemote("ex")
	// THEN
	if len(ids) != 1 {
		t.Error("Expected to return ids of missing exercises")
	}
	if _, ok := refMap["123.fit"]; !ok {
		t.Error("Expected to return a reference to the local exercise")
	}
}
