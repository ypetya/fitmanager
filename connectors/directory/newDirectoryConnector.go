package directory

import "github.com/ypetya/fitmanager/internal"

func NewDirectoryConnector() *DirectoryConnector {

	return &DirectoryConnector{
		io: &internal.EffectiveIO{},
	}
}
