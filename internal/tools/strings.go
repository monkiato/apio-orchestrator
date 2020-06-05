package tools

import "strings"

const (
	alpha        = "abcdefghijklmnopqrstuvwxyz"
	alphanumeric = "abcdefghijklmnopqrstuvwxyz1234567890"
)

func IsAlpha(s string) bool {
	return IsValidFormat(s, alpha)
}

func IsAlphanumeric(s string) bool {
	return IsValidFormat(s, alphanumeric)
}

func IsValidFormat(s string, validCharacters string) bool {
	for _, char := range s {
		if !strings.Contains(validCharacters, strings.ToLower(string(char))) {
			return false
		}
	}
	return true
}

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
