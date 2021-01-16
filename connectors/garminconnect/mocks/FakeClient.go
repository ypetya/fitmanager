package mocks

import (
  "io"
	connect "github.com/abrander/garmin-connect"
)

type FakeClient struct {
	Called     string
	Err        error
	Id         int
	S          int
	C          int
	ActivitiesArr []connect.Activity
  PleasePanicOnActivities bool
  PleasePanicOnExport bool
}

// Implementation of IActivitiesFetcher
func (f *FakeClient) Activities(x string, s int, c int) ([]connect.Activity, error) {
	f.Called, f.S, f.C = "Activities", s, c

  if f.PleasePanicOnActivities {
    panic("Panic on your wish")
  }

	return f.ActivitiesArr, f.Err
}

// Implementation of IActivityExportr
func (f *FakeClient) ExportActivity(id int, w io.Writer, format connect.ActivityFormat) error {
	f.Called, f.Id = "ExportActivity", id

  if f.PleasePanicOnExport {
    panic("Panic on your wish")
  }

	return nil
}
