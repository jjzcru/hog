package hog

import "testing"

func TestGetID(t *testing.T) {
	ids := []string{GetID(), GetID()}

	if len(ids[0]) != len(ids[1]) {
		t.Errorf("the length of both ids should be the same")
	}

	if ids[0] == ids[1] {
		t.Errorf("both id should be different")
	}
}
