package connectors

type IExportCapable interface {
	// Export a single entity to ExportCapable
	// returns the new entity's id
	Export(localDataSourceDir string, localId string) (string, error)
}
