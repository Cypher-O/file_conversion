package utils

import (
	// "strings"
)

// Utility function to check if a slice contains a given string
func Contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
