package utils

import "testing"

func TestRemoveDuplicate(t *testing.T) {
	input := []string{"a", "b", "c", "c"}
	response := RemoveDuplicate(input)
	if len(response) != (len(input) - 1) {
		t.Errorf("the response should have %d items but it has %d instead", len(input)-1, len(response))
	}
}

func TestIsSubstring(t *testing.T) {
	subString := "axc"
	completeString := "111axc111"

	isSubstring, err := IsSubstring(subString, completeString)
	if err != nil {
		t.Error(err)
	}

	if !isSubstring {
		t.Errorf("the substring belong inside the complete string")
	}
}
