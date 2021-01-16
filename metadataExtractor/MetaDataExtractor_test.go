package metadataExtractor

import (
	"testing"
)

func TestExtractShouldNotFailOnFileWithNoSessions(t *testing.T) {
	// WHEN
	a, _, _, _, _, _, _ := Extract("testdata/short.fit")
	// THEN
	if a != "Error" {
		t.Errorf("Should recover with Error type, found %s", a)
	}
}

func TestExtractShouldNotFailOnInvalidFile(t *testing.T) {
	// WHEN
	a, _, _, _, _, _, _ := Extract("testdata/null.fit")
	// THEN
	if a != "Error" {
		t.Errorf("Should recover with Error type, found %s", a)
	}
}

func TestExtractShouldRecoverOnFileErrors(t *testing.T) {
	// WHEN
	a, _, _, _, _, _, _ := Extract("testdata/bad.FIT")
	// THEN
	if a != "Error" {
		t.Errorf("Should recover with Error type, found %s", a)
	}
}

func TestExtract_ShouldExtractMetaDataFrom_Fenix3_File_recorded_outdoor(t *testing.T) {
	// WHEN
	a, d, s, e, ss, b, c := Extract("testdata/1.FIT")
	// THEN
	if a != "Cycling" {
		t.Errorf("Activity is expected to be Cycling, found %s", a)
	}
	if d != "Fenix3" {
		t.Errorf("Device is expected to be Fenix3, found %s", d)
	}
	if s != 1601807125 {
		t.Errorf("Start is expected to be 1601807125 , found %d", s)
	}
	if e != 1601816037 {
		t.Errorf("End is expected to be 1601816037, found %d", e)
	}
	if ss != 2042 {
		t.Errorf("Samples number is expected to be 2042, found %d", ss)
	}
	if len(b) != 6 {
		t.Errorf("Expected to find 6 bands, found %d bands: %s", len(b), b)
	}
	bands := []string{"pos", "hr", "cad", "dist", "cal", "speed"}
	for i, expectedBand := range bands {
		if b[i] != expectedBand {
			t.Errorf("Band %d should be %s, found %s instead.", i, expectedBand, b[i])
		}
	}
	if c != 1601807125 {
		t.Errorf("Created is expected to be 1601807125, found %d", c)
	}
}
func TestExtract_ShouldExtractMetaDataFromZwiftFile_having_no_hr(t *testing.T) {
	// WHEN
	a, d, s, e, ss, b, c := Extract("testdata/zwift.fit")
	// THEN
	if a != "Cycling" {
		t.Errorf("Activity is expected to be Cycling, found %s", a)
	}
	if d != "0" {
		t.Errorf("Device is expected to be 0 for zwift, found %s", d)
	}
	if s != 1604488704 {
		t.Errorf("Start is expected to be 1604488704, found %d", s)
	}
	if e != 1604491600 {
		t.Errorf("End is expected to be 1604491600, found %d", e)
	}
	if ss != 2897 {
		t.Errorf("Samples number is expected to be 2897, found %d", ss)
	}
	if len(b) != 6 {
		t.Errorf("Expected to find 6 bands, found %d bands: %s", len(b), b)
	}
	bands := []string{"pos", "cad", "dist", "pow", "cal", "speed"}
	for i, expectedBand := range bands {
		if b[i] != expectedBand {
			t.Errorf("Band %d should be %s, found %s instead.", i, expectedBand, b[i])
		}
	}

	if c != 1604488681 {
		t.Errorf("Created is expected to be 1604488681, found %d", c)
	}
}
func TestExtract_ShouldExtractMetaDataFromZFenix3File_having_hr(t *testing.T) {
	// WHEN
	a, d, _, _, _, b, _ := Extract("testdata/fenix3.fit")
	// THEN
	if a != "Cycling" {
		t.Errorf("Activity is expected to be Cycling, found %s", a)
	}
	if d != "Fenix3" {
		t.Errorf("Device is expected to be Fenix3 for fenix3, found %s", d)
	}
	if len(b) != 4 {
		t.Errorf("Expected to find 6 bands, found %d bands: %s", len(b), b)
	}
	bands := []string{"hr", "cad", "pow", "cal"}
	for i, expectedBand := range bands {
		if b[i] != expectedBand {
			t.Errorf("Band %d should be %s, found %s instead.", i, expectedBand, b[i])
		}
	}
}
