package connectors

// An union interface for
// * IImportCapable
// * IExportCapable
// * IUpdateCapable
type IConnector interface {
	// Passing the known externalIds (File or Id)
	// the function passes the unknown (new/missing) ids
	// (File or Id)[]
	FetchDiff(knownExternalIds []string) []string
	Import(targetDir string, externalId string) (string, error)
	Export(sourceDir string, filePath string) (string, error)
	Update(localDataSourceDir string, localId string, remoteId string) error
}
