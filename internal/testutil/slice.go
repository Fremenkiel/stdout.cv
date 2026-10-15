package testutil

import (
	"slices"
)

func SliceEql[T comparable](slice1 []T, slice2 []T) bool {
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
