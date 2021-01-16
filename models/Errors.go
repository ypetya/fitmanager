package models

const (
	// Calling method without DataStore.Load
	ErrorDataStoreNotInitialized = Error("DataStore is not initialized")

	ErrorImporterNotFound = Error("Importer not found for a remote type!")
	ErrorExporterNotFound = Error("Exporter not found for a remote type!")
	ErrorRemoteNotFound   = Error("Remote not found!")

	ErrorCouldNotSaveFile = Error("Could not save file")
)
