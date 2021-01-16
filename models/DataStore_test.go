package models

import (
	"os"
	"testing"

	"github.com/ypetya/fitmanager/models/mocks"
)

var fmde mocks.FakeMetadataExtractor
var fio *mocks.FakeIO

// With before/after All:
// remove created tempfiles
// we work with a tempfile on write operations,
// (external dep does not use Reader or Writer interfaces)
func TestMain(m *testing.M) {
	// init mocks
	fmde = mocks.CreateFakeMetadataExtractor()
	fmde.AlwaysNew = true
	fio = &mocks.FakeIO{}

	beforeAll_Export_Suite()

	code := m.Run()

	// teardown
	mocks.FakeIODeleteAllTempfiles()
	os.Exit(code)
}

func Test_DS_ListRemotes_ShouldReturnRemotes(t *testing.T) {
	// GIVEN
	ds := DataStore{
		Remotes: map[string]AnyRemote{
			"garmin": AnyRemote{
				Name:     "garmin",
				LastSync: 0,
				Target:   GarminConnect,
			},
			"ex": AnyRemote{
				Name:     "someexternalfolder",
				LastSync: 0,
				Target:   Directory,
			},
			"l": AnyRemote{
				LastSync: 0,
				Target:   LocalDB,
			},
		},
	}
	// WHEN
	remotes := ds.ListRemotes()
	// THEN
	if len(remotes) != 3 {
		t.Error("Should have 3 remotes")
	}
	if remotes[0].Name != "garmin" && remotes[1].Name != "garmin" && remotes[2].Name != "garmin" {
		t.Error("Should return garmin remote")
	}
}

// When importing a new file into the database it can happen
// that by examining metadata the exercise has been already imported
// When this is the case => A) we could keep both B) keep one of them
// for now we keep the old file on disc, and link the latest
// maybe in the future we will figure out how to extract more information from it
func Test_DS_updateExercise_ShouldTrackExtraFilesInLocalDBFolder(t *testing.T) {
	// GIVEN
	ds := DataStore{
		Exercises: []Exercise{
			Exercise{
				Meta: ActivityMetaInfo{
					Activity: "Cycling",
					Created:  1,
				},
				StoredOn: []AnyRemote{
					AnyRemote{
						LastSync: 1,
						Target:   LocalDB,
						File:     &File{Path: "old.fit"},
					},
					AnyRemote{
						Name:     "garmin",
						LastSync: 0,
						Target:   GarminConnect,
					},
					AnyRemote{
						Name:     "someexternalfolder",
						LastSync: 0,
						Target:   Directory,
					},
				},
			},
		},
	}

	newEx := Exercise{
		Meta: ActivityMetaInfo{
			Activity: "Cycling",
			Created:  1,
		},
		StoredOn: []AnyRemote{
			AnyRemote{
				Target: LocalDB,
				File:   &File{Path: "new.fit"},
			},
		},
	}
	if len(ds.Duplicates) != 0 {
		t.Error("Expected len of Duplicates is 0!")
	}
	// WHEN
	ds.updateExercise(0, &newEx)
	// THEN
	l := len(ds.Duplicates)
	if l != 1 {
		t.Errorf("Expected len of Duplicates is 1, found %d!", l)
	}
	if ds.Duplicates[0].LastSync != 1 {
		t.Errorf("Expected lost file LastSync should remain the original (1), found %d!",
			ds.Duplicates[0].LastSync)
	}
	if ds.Duplicates[0].File.Path != "old.fit" {
		t.Errorf("Expected to save the old.fit file, found %s!",
			ds.Duplicates[0].File.Path)
	}
}

func Test_DataStore_insertExerciseWithRemote_ShouldUpdateTheRemote(t *testing.T) {
	// GIVEN
	mocks.ResetFakeMetadataExtractorCalls()
	// Metadata with prop "Created"=1, 2 existing in testdata !!!
	fmde.FileCreated = map[string]int64{
		"1": 1, "2": 2, // <- these are existing
	}

	ds := DataStore{
		localRootPath:     newFromDir("testdata/remotes-extend-existing"),
		metadataExtractor: fmde,
		Exercises: []Exercise{
			Exercise{
				Meta: ActivityMetaInfo{
					Activity: "Cycling",
					Created:  1,
				},
				StoredOn: []AnyRemote{
					AnyRemote{
						LastSync: 0,
						Target:   LocalDB,
						File:     &File{Path: "any"},
					},
					AnyRemote{
						Name:     "garmin",
						LastSync: 0,
						Target:   GarminConnect,
					},
					AnyRemote{
						Name:     "someexternalfolder",
						LastSync: 0,
						Target:   Directory,
					},
				},
			},
		},
	}

	remote := AnyRemote{
		Name:   "garmin",
		Target: GarminConnect,
	}
	ds.initStruct()
	// WHEN
	ds.insertExerciseWithRemote("AA4C2525.FIT", remote)
	// THEN

	l := len(ds.Exercises[0].StoredOn)
	if l != 3 {
		t.Errorf("Expected to update storedOn entry. Found %d records!", l)
	}
}

func Test_DS_ListRemotes(t *testing.T) {
	// GIVEN
	ds := DataStore{}
	// WHEN
	if !ds.Load("testdata/remotes") {
		t.Errorf("Expected to return Load succeeded")
	}
	// THEN
	r := len(ds.Remotes)
	if r != 3 {
		t.Errorf("Expected to have remotes 3 found %d!", r)
	}
}
func Test_DS_AddRemote(t *testing.T) {
	// GIVEN
	ds := DataStore{}
	// WHEN
	if !ds.Load("testdata/remotes") {
		t.Errorf("Expected to return Load succeeded")
	}
	ds.AddRemote("unknown", AnyRemote{Target: Directory, File: &File{Path: "/tmp"}})
	// THEN
	r := len(ds.Remotes)
	if r != 4 {
		t.Errorf("Expected to have remotes 4 found %d!", r)
	}
}
func Test_DS_AddRemote_of_GarminConnect(t *testing.T) {
	// GIVEN
	ds := DataStore{}
	// WHEN
	if !ds.Load("testdata/remotes") {
		t.Errorf("Expected to return Load succeeded")
	}
	ds.AddRemote("garmin2", AnyRemote{Target: GarminConnect})
	// THEN
	r := len(ds.Remotes)
	if r != 4 {
		t.Errorf("Expected to have remotes 4 found %d!", r)
	}
}

func Test_DS_AddRemote_ShouldFailWhenNoSuchDirectory(t *testing.T) {
	// GIVEN
	ds := DataStore{}
	// WHEN
	if !ds.Load("testdata/remotes") {
		t.Errorf("Expected to return Load succeeded")
	}
	ds.AddRemote("unknown", AnyRemote{Target: Directory, File: &File{Path: "tmp2"}})
	// THEN
	r := len(ds.Remotes)
	if r != 3 {
		t.Errorf("Expected to have remotes 3 found %d!", r)
	}
}

func Test_DS_DelRemote(t *testing.T) {
	// GIVEN
	ds := DataStore{}
	// WHEN
	if !ds.Load("testdata/remotes") {
		t.Errorf("Expected to return Load succeeded")
	}
	ds.DelRemote("zwift")
	// THEN
	r := len(ds.Remotes)
	if r != 2 {
		t.Errorf("Expected to have remotes 2 found %d!", r)
	}
}

func Test_DS_getRemoteRefs(t *testing.T) {
	// GIVEN
	ds := DataStore{}
	if !ds.Load("testdata/remotes") {
		t.Errorf("Expected to return Load succeeded")
	}
	// WHEN
	ids := ds.collectStoredRefsForRemote("garmin")
	// THEN

	if len(ids) != 2 || ids[0] != "1" || ids[1] != "2" {
		t.Error("Expected to return ids(1,2) for exercises")
	}
}

func Test_DataStore_Should_return_ErrorNotInitialized_when_save_is_called_without_load(t *testing.T) {
	// GIVEN
	fio := &mocks.FakeIO{}
	ds := DataStore{
		localRootPath: newFromDir("testdata/valid"),
		io:            fio,
	}
	// WHEN
	err := ds.Save()
	// THEN
	if err != ErrorDataStoreNotInitialized {
		t.Errorf("Expected to return ErrorDataStoreNotInitialized")
	}
}

func Test_DS_Save_ShouldSerialize(t *testing.T) {
	// GIVEN
	fio := &mocks.FakeIO{}
	ds := DataStore{
		io: fio,
	}
	// WHEN
	ds.Load("testdata/valid")
	err := ds.Save()
	// THEN
	if err != nil {
		t.Errorf("Expected to return false, as not initialized!")
	}
}
func Test_DS_LoadValidDataSet_ShouldDeserialize(t *testing.T) {
	// GIVEN
	ds := DataStore{}
	// WHEN
	succ := ds.Load("testdata/valid")
	// THEN
	if succ == false {
		t.Errorf("Expected to return Load succeeded")
	}
	if len(ds.Exercises) != 1 {
		t.Errorf("Expected to deserialize exercises array!")
	}
}

func Test_DS_LoadEmptyDataSet_ShouldDeserialize(t *testing.T) {
	// GIVEN
	ds := DataStore{}
	// WHEN
	ds.Load("testdata/empty")
	// THEN
	if ds.Version != 1 {
		t.Errorf("Expected to deserialize version!")
	}
	if len(ds.Exercises) != 0 {
		t.Errorf("Expected to deserialize empty exercises array!")
	}
}

func Test_DataStore_when_Load_is_called_on_an_empty_dir_without_a_dataSet(t *testing.T) {
	// GIVEN
	ds := DataStore{}
	// WHEN
	succ := ds.Load("testdata/empty_folder")
	// THEN
	if len(ds.Exercises) != 0 {
		t.Errorf("Expected exercises num is 0 (found %d)", len(ds.Exercises))
	}
	t1 := ds.lastSync
	if t1 == 0 {
		t.Errorf("Expected to have lastSync!")
	}
	if succ {
		t.Errorf("Load expected to return false!")
	}
}

func Test_DataStore_when_Save_is_called_on_an_empty_dir_without_a_dataSet(t *testing.T) {
	// GIVEN
	fio := &mocks.FakeIO{}
	ds := DataStore{
		io: fio,
	}
	// WHEN
	ds.Load("testdata/empty_folder")
	ds.AddRemote("garmin", AnyRemote{
		Target: GarminConnect,
		Name:   "garmin",
	})
	err := ds.Save()

	// THEN
	if err != nil {
		t.Errorf("Expected to have Save succeeded")
	}
}

func Test_DS_LoadEmptyDataSet_ShouldSync(t *testing.T) {
	// GIVEN
	ds := DataStore{}
	// WHEN
	succ := ds.Load("testdata/empty")
	// THEN
	if succ == false {
		t.Errorf("Expected to return Load succeeded")
	}

	t1 := ds.lastSync
	if t1 == 0 {
		t.Errorf("Expected to have lastSync!")
	}
	t2 := ds.localRootPath.LastSeen
	if t1 != t2 {
		t.Errorf("Expected to update lastSync timestamp to data store path access-time! Found different timestamps %d vs %d", t1, t2)
	}
}

func Test_DS_LoadDirWithoutDir(t *testing.T) {
	// GIVEN
	ds := DataStore{}
	// WHEN
	ret := ds.Load("testdata/file.txt")
	// THEN
	if len(ds.Exercises) != 0 {
		t.Errorf("Expected exercises num is 0 (found %d)", len(ds.Exercises))
	}
	if ret != false {
		t.Errorf("Expected to return Load false!")
	}
}

func Test_DS_Save_Should_Order_Exercises(t *testing.T) {
	// GIVEN
	fio := &mocks.FakeIO{}
	ds := DataStore{
		io: fio,
	}

	succ := ds.Load("testdata/order")
	if succ == false {
		t.Errorf("Expected to return Load succeeded")
	}
	created1, created2 := ds.Exercises[0].Meta.Created,
		ds.Exercises[1].Meta.Created
	if created1 != 2 && created2 != 1 {
		t.Errorf("Input should be unordered")
	}
	// WHEN
	err := ds.Save()

	// THEN
	if err != nil {
		t.Errorf("Expected to have Save succeeded")
	}
	if len(ds.Exercises) != 2 {
		t.Errorf("Expected to still have 2 exercises!")
	}
	created1, created2 = ds.Exercises[0].Meta.Created,
		ds.Exercises[1].Meta.Created
	if created1 != 1 && created2 != 2 {
		t.Errorf("Expected to have exercises ordered!")
	}
}

func Test_DS_Save_Should_Set_OverlapIds_On_Exercises(t *testing.T) {
	// GIVEN
	fio := &mocks.FakeIO{}
	ds := DataStore{
		io: fio,
	}

	succ := ds.Load("testdata/overlap")
	if succ == false {
		t.Errorf("Expected to return Load succeeded")
	}
	ovr1, ovr2 := len(ds.Exercises[0].OverlappingIds),
		len(ds.Exercises[1].OverlappingIds)
	if ovr1 != 0 && ovr2 != 0 {
		t.Errorf("Input should have no overlapping ids")
	}
	// WHEN
	err := ds.Save()

	// THEN
	if err != nil {
		t.Errorf("Expected to have Save succeeded")
	}
	if len(ds.Exercises) != 2 {
		t.Errorf("Expected to still have 2 exercises!")
	}
	ovr1, ovr2 = len(ds.Exercises[0].OverlappingIds),
		len(ds.Exercises[1].OverlappingIds)
	if ovr1 != 1 || ovr2 != 1 {
		t.Errorf("Expected to have overlaps registered in exercises found #[%d,%d]", ovr1, ovr2)
	}
}
