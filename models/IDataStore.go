package models

type IDataStore interface {
	Import(string) error
	Export(string) error
	AddRemote(string, AnyRemote)
	DelRemote(string)
	ListRemotes() []AnyRemote
	Load(path string) bool
	Save() error
	GetExercises() []Exercise
	Filter(f IFilter) []Exercise
	NewRemoteFilter() IRemoteFilter
}
