package utils

import (
	"sort"
)

// MergeAndRemoveDuplicates merges two slices of strings and removes duplicates
func MergeAndRemoveDuplicates(first []string, second []string) []string {
	// Create a map to track unique values
	uniqueMap := make(map[string]bool)

	// Add all values from first to the map
	for _, value := range first {
		uniqueMap[value] = true
	}

	// Add values from second to the map if they are not already present
	for _, value := range second {
		if _, exists := uniqueMap[value]; !exists {
			uniqueMap[value] = true
		}
	}

	// Convert the unique map keys to a new slice
	merged := make([]string, 0, len(uniqueMap))
	for key := range uniqueMap {
		merged = append(merged, key)
	}

	// Sort the slice
	sort.Strings(merged)

	return merged
}
