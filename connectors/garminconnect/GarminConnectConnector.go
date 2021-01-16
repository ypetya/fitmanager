package garminconnect

import (
	"github.com/ypetya/fitmanager/internal"
)

// implements IConnector(IImportCapable, IExportCapable)
type GarminConnectConnector struct {
	fetcher       IActivitiesFetcher
	exporter      IActivityExporter
	authenticator IAuthenticator
	authenticated bool
	io            internal.IFileCreator
}
