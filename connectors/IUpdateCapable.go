package connectors

type IUpdateCapable interface {
	// Passing the known externalIds (File or Id)
	// the function passes the unknown (new/missing) ids
	// (File or Id)[]
	FetchDiff(knownExternalIds []string) []string
	// Update function
	Update(localDataSourceDir string, localId string, remoteId string) error
}
