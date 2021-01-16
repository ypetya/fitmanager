package garminconnect

import (
	"testing"

	connect "github.com/abrander/garmin-connect"

	"github.com/ypetya/fitmanager/connectors/garminconnect/mocks"
	"github.com/ypetya/fitmanager/internal"
)

func TestFetchActivitiesCallsClient(t *testing.T) {
	// GIVEN
	fc := &mocks.FakeClient{ActivitiesArr: []connect.Activity{}}
	// WHEN
	act := fetchActivities(fc, 12, 24)
	// THEN
	if fc.Called != "Activities" {
		t.Errorf("Expected to call client.Activities")
	}
	if fc.S != 12 || fc.C != 24 {
		t.Errorf("Expected to call with offset 12 (found %d) and count 24 (found %d)", fc.S, fc.C)
	}
	if len(act) != 0 {
		t.Errorf("Expected to return an empty array")
	}
}

func TestFetchActivitiesShouldRecover(t *testing.T) {
	// GIVEN
	fc := &mocks.FakeClient{ActivitiesArr: []connect.Activity{},
		PleasePanicOnActivities: true}
	// WHEN
	act := fetchActivities(fc, 12, 24)
	// THEN
	if fc.Called != "Activities" {
		t.Errorf("Expected to call client.Activities")
	}
	if fc.S != 12 || fc.C != 24 {
		t.Errorf("Expected to call with offset 12 (found %d) and count 24 (found %d)", fc.S, fc.C)
	}
	if len(act) != 0 {
		t.Errorf("Expected to return an empty array")
	}
}

func TestDownloadShouldRecover(t *testing.T) {
	// GIVEN
	fc := &mocks.FakeClient{
		PleasePanicOnExport: true,
	}
	fio := &internal.FakeIO{}
	// WHEN
	str := downloadActivity(fc, fio, "dir/", 12)
	// THEN
	if fc.Called != "ExportActivity" || fc.Id != 12 {
		t.Errorf("Expected to call ExportActivity on client id 12 (found %d)", fc.Id)
	}
	if fio.Called != "createFile" || fio.Name != "dir/12.fit" {
		t.Errorf("Expected to call createFileForWriting with name 12.fit, found(%s)", fio.Name)
	}
	if str != "12.fit" {
		t.Errorf("Expected to return the filename")
	}
}

// when a file already downloaded but not in db, it should skip
// the export call
func TestDownloadShouldSkipExistingFiles(t *testing.T) {
	// GIVEN
	fc := &mocks.FakeClient{}
	fio := &internal.FakeIO{}
	// WHEN
	str := downloadActivity(fc, fio, "testdata/", 2)
	// THEN
	if fc.Called == "ExportActivity" {
		t.Errorf("Expected _NOT_ to call ExportActivity on client")
	}
	if fio.Called == "createFile" {
		t.Errorf("Expected _NOT_ to call createFileForWriting")
	}
	if str != "2.fit" {
		t.Errorf("Expected to return the filename")
	}
}

func TestDownloadActivityCalls(t *testing.T) {
	// GIVEN
	fc := &mocks.FakeClient{}
	fio := &internal.FakeIO{}
	// WHEN
	str := downloadActivity(fc, fio, "dir/", 12)
	// THEN
	if fc.Called != "ExportActivity" || fc.Id != 12 {
		t.Errorf("Expected to call ExportActivity on client id 12 (found %d)", fc.Id)
	}
	if fio.Called != "createFile" || fio.Name != "dir/12.fit" {
		t.Errorf("Expected to call createFileForWriting with name 12.fit, found(%s)", fio.Name)
	}
	if str != "12.fit" {
		t.Errorf("Expected to return the filename")
	}
}
