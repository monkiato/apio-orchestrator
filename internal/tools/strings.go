package tools

import "strings"

const (
	alpha        = "abcdefghijklmnopqrstuvwxyz"
	alphanumeric = "abcdefghijklmnopqrstuvwxyz1234567890"
)

//IsAlpha confirm if a string is using alphabetic characters only
func IsAlpha(s string) bool {
	return IsValidFormat(s, alpha)
}

//IsAlphanumeric confirm if a string is using alphanumeric characters only
func IsAlphanumeric(s string) bool {
	return IsValidFormat(s, alphanumeric)
}

//IsValidFormat confirm if a string is using the specified list of acceptable characters
func IsValidFormat(s string, validCharacters string) bool {
	for _, char := range s {
		if !strings.Contains(validCharacters, strings.ToLower(string(char))) {
			return false
		}
	}
	return true
}

//Find a specific string inside a slice of strings
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
