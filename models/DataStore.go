package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"time"

	"github.com/ypetya/fitmanager/connectors"
	directoryConnector "github.com/ypetya/fitmanager/connectors/directory"
	"github.com/ypetya/fitmanager/internal"
	"github.com/ypetya/fitmanager/metadataExtractor"
)

var effectiveIO = internal.EffectiveIO{}

// This is the main entrypoint of the database - exercises stored locally
type DataStore struct {
	// serialization version id stored in ds.json - equals to CMD application version
	Version int64 `json:"version"`
	// exercises with metadata stored in ds.json
	Exercises []Exercise `json:"exercises,omitempty"`
	// remotes config stored in ds.json
	Remotes map[string]AnyRemote `json:"remotes"`
	// store a list of files which have been overwritten by an import
	Duplicates []AnyRemote `json:"duplicates,omitempty"`
	// time-stamp saved / loaded
	lastSync int64
	// where the local datastore exists
	localRootPath File
	// struct initialized - internal maps, effective file IO attached
	initialized bool
	// internal ref map for exercise lookup
	refMap      map[string]*Exercise
	RefIndexMax int `json:"ix,omitempty"`

	// --- Abstract producers ---

	// write file operations to mock disk-writes
	io         internal.IFileCreator
	connectors map[RemoteType]connectors.IConnector
	// produce meta-data from fit files
	metadataExtractor metadataExtractor.IMetadataExtractor
	// definition of enhancers
	enhancers []Enhancer
}

// sets the internal Exercises RefIndex
func (db *DataStore) buildRefMap() {
	for ix, _ := range db.Exercises {
		ptr := &db.Exercises[ix]
		db.createRef(ptr)
	}
}

func (db *DataStore) createRef(e *Exercise) {
	if ref, err := e.GetLocalRef(); err == nil {
		db.refMap[ref] = e
		if e.RefIndex == 0 {
			db.RefIndexMax = db.RefIndexMax + 1
			e.RefIndex = db.RefIndexMax
		}
	}
}

func (db *DataStore) Export(remoteName string) error {
	// init
	remoteTemplate, found := db.Remotes[remoteName]
	if !found {
		return ErrorRemoteNotFound
	}
	exporter, found := db.connectors[remoteTemplate.GetType()]
	if !found {
		return ErrorImporterNotFound
	}
	if len(db.Exercises) == 0 {
		return nil
	}
	// setup exporter
	if remoteTemplate.GetType() == Directory {
		if dirExporter, ok := exporter.(directoryConnector.IDirectoryExporter); ok {
			dirExporter.SetTarget(remoteTemplate.File.AsDir())
		}
	}
	// enhance
	if err := db.enhance(); err != nil {
		return err
	}
	// Export
	db.exportNewExercises(exporter, &remoteTemplate)
	db.exportRemovedExercises(exporter, remoteName)

	return db.Save()
}

// New exercises: not having the remote, works on overlapping as well
func (db *DataStore) exportNewExercises(
	exporter connectors.IExportCapable,
	remoteTemplate *AnyRemote) {
	dir := db.localRootPath.AsDir()
	remoteName := remoteTemplate.Name
	for ix, _ := range db.Exercises {
		ep := &db.Exercises[ix]
		if db.isExerciseValidForUpload(ep, remoteName) {
			fileName, _ := ep.GetLocalRef()
			if newRemoteId, err := exporter.Export(dir, fileName); err == nil {
				remote := remoteTemplate.Fill(newRemoteId)
				db.addExerciseRemote(ep, remote)
			}
		}
	}
}

func (db *DataStore) exportRemovedExercises(exporter connectors.IUpdateCapable, remoteName string) {
	dir := db.localRootPath.AsDir()
	existingIds, refMap := db.collectExercisesForRemote(remoteName)
	newIds := exporter.FetchDiff(existingIds)
	for _, remoteId := range newIds {
		// When not found in refMap provided by collectExercises means
		// FetchDiff returned this because it is a normal export not an update
		// ~ diff === (a)missing on remote UNION (b)missing on local
		// a -> shouldUpdate
		// b -> shouldExport !
		if ex, shouldUpdate := refMap[remoteId]; shouldUpdate {
			fileName, _ := ex.GetLocalRef()
			if err := exporter.Update(dir, fileName, remoteId); err != nil {
				fmt.Printf("Could not export %s %s\n", remoteId, err)
			}
		}
	}
}

func (db *DataStore) enhance() error {
	if len(db.Exercises) == 0 {
		return nil
	}
	dir := db.localRootPath.AsDir()
	for _, enhancer := range db.enhancers {
		exercises := db.Filter(enhancer.Filter)
		for _, exercise := range exercises {
			if toEnhance, err := exercise.GetLocalRef(); err == nil {
				// FIXME make it work for multiple overlaps
				with := dir + exercise.OverlappingIds[0]
				// TODO introduce filename formatter -> enhancer
				target := toEnhance + "." + enhancer.Name + ".fit"
				enhancer.Function.Enhance(dir+target, dir+toEnhance, with)
				db.insertExerciseWithLocalRef(target)
			}
		}
	}

	// save sync + mark overlap
	return db.Save()
}

func (db *DataStore) GetExercises() []Exercise {
	return db.Exercises
}

// @returns a RemoteFilter
func (db *DataStore) NewRemoteFilter() IRemoteFilter {
	return &RemoteFilter{}
}

// f IFilter - accepts any filter which has a Filter a method
func (db *DataStore) Filter(f IFilter) []Exercise {
	return f.Filter(db.Exercises)
}

// Imports from a specific remoteName
// a remote name is a string, which points to a specific instance of
// a specific typed remote
// For example you can have multiple remotes for different GarminConnect
// accounts. This is a future goal to implement.
func (db *DataStore) Import(remoteName string) error {
	// init
	dir := db.localRootPath.AsDir()
	remoteTemplate, found := db.Remotes[remoteName]
	if !found {
		return ErrorRemoteNotFound
	}
	importer, found := db.connectors[remoteTemplate.GetType()]
	if !found {
		return ErrorImporterNotFound
	}
	// setup importer
	if remoteTemplate.GetType() == Directory {
		if dirImporter, ok := importer.(directoryConnector.IDirectoryImporter); ok {
			dirImporter.SetSource(remoteTemplate.File.AsDir())
		}
	}
	// import
	existingIds := db.collectStoredRefsForRemote(remoteName)
	newIds := importer.FetchDiff(existingIds)
	for _, id := range newIds {
		if fileName, err := importer.Import(dir, id); err == nil {
			remote := remoteTemplate.Fill(id)
			db.insertExerciseWithRemote(fileName, remote)
		}
	}

	// save sync
	remoteTemplate.LastSync = db.ts()
	db.Remotes[remoteName] = remoteTemplate
	return db.Save()
}

func (db *DataStore) addExerciseRemote(ex *Exercise, newRemote AnyRemote) {
	for _, remote := range ex.StoredOn {
		if newRemote.equals(remote) {
			return
		}
	}
	ex.StoredOn = append(ex.StoredOn, newRemote)
}

// Inserts an exercise given the localRef
func (db *DataStore) insertExerciseWithLocalRef(ref string) bool {
	fileName := ref
	dir := db.localRootPath.AsDir()

	activity, device, start, end, samples, bands, created :=
		db.metadataExtractor.Extract(dir + fileName)

	newEx := Exercise{
		Meta: ActivityMetaInfo{
			Activity: activity,
			Device:   device,
			Start:    start,
			End:      end,
			Samples:  samples,
			Bands:    bands,
			Created:  created,
		},
		StoredOn: []AnyRemote{
			AnyRemote{
				Target: LocalDB,
				File: &File{
					Path:     fileName,
					LastSeen: db.ts(),
				},
			},
		},
	}
	db.createRef(&newEx)

	for i, e := range db.Exercises {
		if e.equals(newEx) {
			db.updateExercise(i, &newEx)
			return false
		}
	}

	db.Exercises = append(db.Exercises, newEx)

	return true
}

// extend db puts a new exercise record into the database might it be
// an insert or an update record
// @returns true when a new record has been inserted
func (db *DataStore) insertExerciseWithRemote(
	ref string,
	updatedRemote AnyRemote) bool {
	fileName := ref
	dir := db.localRootPath.AsDir()

	activity, device, start, end, samples, bands, created :=
		db.metadataExtractor.Extract(dir + fileName)

	newEx := Exercise{
		Meta: ActivityMetaInfo{
			Activity: activity,
			Device:   device,
			Start:    start,
			End:      end,
			Samples:  samples,
			Bands:    bands,
			Created:  created,
		},
		StoredOn: []AnyRemote{
			AnyRemote{
				Target: LocalDB,
				File: &File{
					Path:     fileName,
					LastSeen: db.ts(),
				},
			},
			updatedRemote,
		},
	}
	db.createRef(&newEx)

	for i, e := range db.Exercises {
		if e.equals(newEx) {
			db.updateExercise(i, &newEx)
			return false
		}
	}

	db.Exercises = append(db.Exercises, newEx)

	return true
}

// Updates an exercise in the database with merging the remotes
// if a remote has the same store, the newer reference gets stored
// and the old one gets moved to the duplicates records.
// These can be later garbage collected or reprocessed when meta-data
// extracting mechanism gets improved
func (db *DataStore) updateExercise(ix int, newEx *Exercise) {
	// existing records moved to the new exercise
original:
	for i, oldRemote := range db.Exercises[ix].StoredOn {
		for _, newRemote := range newEx.StoredOn {
			if newRemote.equals(oldRemote) {
				db.Duplicates = append(db.Duplicates, db.Exercises[ix].StoredOn[i])
				obsoleteRef := oldRemote.GetRef()
				delete(db.refMap, obsoleteRef)
				continue original
			}
		}
		newEx.StoredOn = append(newEx.StoredOn, db.Exercises[ix].StoredOn[i])
	}

	// update
	db.Exercises[ix] = *newEx
}

func (db *DataStore) collectStoredRefsForRemote(name string) []string {
	var ret []string

	for _, e := range db.Exercises {
		ref, err := e.GetRemoteRef(name)
		if err == nil {
			ret = append(ret, ref)
		}
	}
	return ret
}

func (db *DataStore) collectExercisesForRemote(name string) (ret []string, refMap map[string]*Exercise) {
	refMap = make(map[string]*Exercise)
	for ix, e := range db.Exercises {
		ref, err := e.GetRemoteRef(name)
		if err == nil {
			ret = append(ret, ref)
			refMap[ref] = &db.Exercises[ix]
		}
	}
	return ret, refMap
}

func (db *DataStore) AddRemote(name string, r AnyRemote) {
	switch {
	case r.Target == GarminConnect:
		db.Remotes[name] = r
	case r.Target == Directory:
		if db.io.FileExists(r.File.AsDir()) {
			db.Remotes[name] = r
		}
	}
}
func (db *DataStore) DelRemote(name string) {
	delete(db.Remotes, name)
}
func (db *DataStore) ListRemotes() []AnyRemote {
	ret := []AnyRemote{}
	for key, _ := range db.Remotes {
		ret = append(ret, db.Remotes[key])
	}
	return ret
}

// @returns success true when directory valid and ds.json loaded
func (db *DataStore) Load(path string) bool {
	db.localRootPath = newFromDir(path)
	db.initStruct()
	dirExist := db.markSync(LocalDB)

	if dirExist {
		err := db.loadDataSetMetaData()
		if err != nil {
			return false
		}
	}

	return dirExist
}

func (db *DataStore) ts() int64 {
	return time.Now().Unix()
}

func (db *DataStore) Save() (named error) {
	if !db.initialized {
		return ErrorDataStoreNotInitialized
	}

	db.orderExercises()
	db.markOverlappingExercises()

	db.localRootPath.LastSeen = db.ts()
	err := db.saveDataSetMetaData()
	if err == nil {
		db.markSync(LocalDB)
		return nil
	}
	return err
}

func (db *DataStore) orderExercises() {
	sort.Sort(db)
}

func (db *DataStore) fileName() string {
	return db.localRootPath.Path + "/ds.json"
}

// This method does the physical loading
// deserializing metadata
func (db *DataStore) loadDataSetMetaData() error {
	fileName := db.fileName()
	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, db)

	if err != nil {
		return err
	}
	db.buildRefMap()

	return nil
}

// This method does the physical writing
// serializing metadata
func (db *DataStore) saveDataSetMetaData() error {
	fileName := db.fileName()
	bytes, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return err
	}
	f := db.io.CreateFileForWriting(fileName)
	_, err = f.Write(bytes)
	defer f.Close()
	return err
}

func (db *DataStore) initStruct() {
	db.Remotes = make(map[string]AnyRemote)
	db.refMap = make(map[string]*Exercise)

	// Attach effective file writer
	if db.io == nil {
		db.io = effectiveIO
	}

	db.initialized = true
}

func (db *DataStore) markSync(t RemoteType) bool {
	db.lastSync = db.localRootPath.LastSeen

	return db.localRootPath.LastSeen > 0
}
