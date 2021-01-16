package directory

import (
	"os"
	"path/filepath"

	"github.com/ypetya/fitmanager/internal"
)

type DirectoryConnector struct {
	io  internal.IFileCreator
	dir string
}

// Returns the union of not found on remote and not found in db
// This function both works...
// for import: it returns the not found in db (new files)
// for export: it returns the not found in directory (new files)
func (di *DirectoryConnector) FetchDiff(knownExternalIds []string) []string {
	found := []string{}
	f, err := os.Open(di.dir)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	fileList, err := f.Readdir(-1)
	if err != nil {
		panic(err)
	}
	fileNames := []string{}
	for _, fileInfo := range fileList {
		if !fileInfo.IsDir() && fileInfo.Mode().IsRegular() {
			candidate := fileInfo.Name()
			ext := filepath.Ext(candidate)
			if ext == ".FIT" || ext == ".fit" {
				fileNames = append(fileNames, candidate)
			}
		}
	}
outer:
	for _, candidate := range fileNames {
		for _, existing := range knownExternalIds {
			if existing == candidate {
				continue outer
			}
		}
		found = append(found, candidate)
	}
outer2:
	for _, known := range knownExternalIds {
		for _, fileName := range fileNames {
			if fileName == known {
				continue outer2
			}
		}
		found = append(found, known)
	}
	return found
}

func (di *DirectoryConnector) Import(targetDir string, externalId string) (string, error) {
	src := di.dir + externalId
	dst := targetDir + externalId
	di.io.Copy(src, dst)
	return externalId, nil
}

func (di *DirectoryConnector) Export(sourceDir string, id string) (string, error) {
	src := sourceDir + id
	dst := di.dir + id
	di.io.Copy(src, dst)
	return id, nil
}

func (di *DirectoryConnector) Update(localDataSourceDir string, localId string, remoteId string) error {
	src := localDataSourceDir + localId
	dst := di.dir + remoteId
	di.io.Copy(src, dst)
	return nil
}

func (di *DirectoryConnector) SetSource(dir string) {
	di.dir = dir
}

func (di *DirectoryConnector) SetTarget(dir string) {
	di.dir = dir
}
