package utils

import "slices"

func Unique[T comparable](s []T) []T {
	result := make([]T, 0)
	for _, elem := range s {
	  if !slices.Contains(result, elem) {
		  result = append(result, elem)
	  }
	}

	return result
}
