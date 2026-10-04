package testutil

import "slices"

func SliceEql(slice1 []string, slice2 []string) bool {
	if s1Len := len(slice1); len(slice2) != s1Len {
		return false
	}

	for _, s := range slice1 {
		if !slices.Contains(slice2, s) {
			return false
		}
	}
	return true
}
