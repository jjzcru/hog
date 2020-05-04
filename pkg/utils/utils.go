package utils

import (
	"fmt"
	"regexp"
)

// RemoveDuplicate takes a slice of string as input and remove the duplicates
func RemoveDuplicate(input []string) []string {
	var response []string
	if input == nil {
		input = []string{}
	}

	responseMap := map[string]bool{}

	for _, value := range input {
		responseMap[value] = true
	}

	for k := range responseMap {
		response = append(response, k)
	}

	return response
}

// IsSubstring check if a string is a substring of another string
func IsSubstring(substring string, completeString string) (bool, error) {
	reg, err := regexp.Compile(fmt.Sprintf("(.)*%s(.)*", substring))
	if err != nil {
		return false, err
	}
	return reg.MatchString(completeString), nil
}
