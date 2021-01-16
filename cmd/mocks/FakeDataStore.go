package mocks

import "github.com/ypetya/fitmanager/models"

type FakeDataStore struct {
	Calls        []string
	Passed       []string
	LoadReturn   bool
	Exercises    []models.Exercise
	FilterPassed *models.IFilter
}

func (f *FakeDataStore) GetExercises() []models.Exercise {
	f.Calls = append(f.Calls, "GetExercises")
	return f.Exercises
}

func (f *FakeDataStore) NewRemoteFilter() models.IRemoteFilter {
	return &FakeFilter{}
}

func (f *FakeDataStore) Filter(filter models.IFilter) []models.Exercise {
	f.Calls = append(f.Calls, "Filter")
	f.FilterPassed = &filter
	return f.Exercises
}

func (f *FakeDataStore) Import(name string) error {
	f.Calls = append(f.Calls, "Import")
	f.Passed = append(f.Passed, name)
	return nil
}
func (f *FakeDataStore) Export(name string) error {
	f.Calls = append(f.Calls, "Export")
	f.Passed = append(f.Passed, name)
	return nil
}
func (f *FakeDataStore) DelRemote(name string) {
	f.Calls = append(f.Calls, "DelRemote")
	f.Passed = append(f.Passed, name)
}
func (f *FakeDataStore) AddRemote(name string, _ models.AnyRemote) {
	f.Calls = append(f.Calls, "AddRemote")
	f.Passed = append(f.Passed, name)
}
func (f *FakeDataStore) ListRemotes() []models.AnyRemote {
	f.Calls = append(f.Calls, "ListRemotes")
	return []models.AnyRemote{}
}
func (f *FakeDataStore) Load(dir string) bool {
	f.Calls = append(f.Calls, "Load")
	return f.LoadReturn
}
func (f *FakeDataStore) Save() error {
	f.Calls = append(f.Calls, "Save")
	return nil
}
func (f *FakeDataStore) LastSeen() int64 {
	f.Calls = append(f.Calls, "LastSeen")
	return 0
}
