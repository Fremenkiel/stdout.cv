package testutil

import (
	"slices"
)

func SliceEql(slice1 []string, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	for _, s := range slice1 {
		if !slices.Contains(slice2, s) {
			return false
		}
	}
	return true
}
