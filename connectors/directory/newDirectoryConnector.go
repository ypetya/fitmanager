package directory

import "../../internal"

func NewDirectoryConnector() *DirectoryConnector {

	return &DirectoryConnector{
		io: &internal.EffectiveIO{},
	}
}
