package connectors

type IImportCapable interface {
	// Passing the known externalIds (File or Id)
	// the function passes the unknown (new/missing) ids
	// (File or Id)[]
	FetchDiff(knownExternalIds []string) []string
	// Import a single external entity to targetDir
	// returns the new file's name
	Import(targetDir string, externalId string) (string, error)
}
