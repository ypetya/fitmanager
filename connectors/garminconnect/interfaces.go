package garminconnect

import (
	connect "github.com/abrander/garmin-connect"
	"io"
)

// to mock out connect.Activities call
type IActivitiesFetcher interface {
	Activities(string, int, int) ([]connect.Activity, error)
}

// to mock out connect.ExportActivity call
type IActivityExporter interface {
	ExportActivity(int, io.Writer, connect.ActivityFormat) error
}

type IAuthenticator interface {
	Authenticate() error
	SetOptions(options ...connect.Option)
}
