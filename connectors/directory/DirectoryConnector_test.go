package directory

import (
	"testing"

	"../../internal"
)

func Test_DirectoryConnector_FetchDiff_ShouldReturn_FitFiles_Not_Imported(t *testing.T) {
	// GIVEN
	di := DirectoryConnector{}
	di.SetSource("testdata/")

	known := []string{"1.FIT"}
	// WHEN
	found := di.FetchDiff(known)
	// THEN
	if len(found) != 2 {
		t.Errorf("Return only newly found *.fit and *.FIT files, found %d", len(found))
		t.Error(found)
	}
}

func Test_DirectoryConnector_Export_Should_call_copy_on_concatenated_path(t *testing.T) {
	// GIVEN
	fio := internal.FakeIO{}
	di := DirectoryConnector{
		io: &fio, dir: "target/",
	}
	// WHEN
	fn, err := di.Export("ds/", "2.fit")
	// THEN
	if fio.Called != "Copy" {
		t.Error("Should call io.Copy!")
	}
	if fio.Dst != "target/2.fit" {
		t.Error("Should call with dest 'target/2.fit'")
	}
	if fio.Name != "ds/2.fit" {
		t.Error("Should call with source 'ds/2.fit'")
	}
	if err != nil {
		t.Error("Should return err==nil!")
	}
	if fn != "2.fit" {
		t.Error("Should return fileName without dir '2.fit'")
	}
}

func Test_DirectoryConnector_Import_Should_Call_Copy_With_ConcatenatedPath(t *testing.T) {
	// GIVEN
	fio := internal.FakeIO{}
	di := DirectoryConnector{
		io: &fio, dir: "testdata/",
	}
	// WHEN
	fn, err := di.Import("target/", "2.fit")
	// THEN
	if fio.Called != "Copy" {
		t.Error("Should call io.Copy!")
	}
	if fio.Dst != "target/2.fit" {
		t.Error("Should call with dest 'target/2.fit'")
	}
	if fio.Name != "testdata/2.fit" {
		t.Error("Should call with source 'testdata/2.fit'")
	}
	if err != nil {
		t.Error("Should return err==nil!")
	}
	if fn != "2.fit" {
		t.Error("Should return fileName without dir '2.fit'")
	}
}
