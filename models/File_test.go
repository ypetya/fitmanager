package models

import "testing"

func TestFileWithPathShouldHaveAName(t *testing.T) {
	// GIVEN
	f := File{Path: "example.fit"}
	// WHEN
	e := f.HasName()
	// THEN
	if !e {
		t.Error("Expected HasName to be true")
	}
}

func TestEmptyFileShouldNotHaveAName(t *testing.T) {
	// GIVEN
	f := File{}
	// WHEN
	e := f.HasName()
	// THEN
	if e {
		t.Error("Expected HasName to be false")
	}
}

func TestAsDirEnsuresLastChar(t *testing.T) {
	// GIVEN
	filesIn := []string{"a", "b/", "", "o/./"}
	expected := []string{"a/", "b/", "", "o/./"}
	//WHEN
	for i, _ := range filesIn {
		fin := File{Path: filesIn[i]}
		if fin.AsDir() != expected[i] {
			t.Errorf("Invalid format for %s should be %s", filesIn[i], expected[i])
		}
	}
}
