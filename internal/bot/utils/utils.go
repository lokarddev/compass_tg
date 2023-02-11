package utils

import "sort"

func StrInSlice(slice []string, need string) bool {
	for _, item := range slice {
		if item == need {
			return true
		}
	}

	return false
}

func IntInSlice(slice []int, need int) bool {
	for _, item := range slice {
		if item == need {
			return true
		}
	}

	return false
}

func RemoveDuplicatesAndSortStr(initialSlice []string) []string {
	processedSlice := RemoveDuplicatesF(initialSlice)
	sort.Strings(processedSlice)

	return processedSlice
}

func RemoveDuplicatesF(s []string) []string {
	uniqueMap := make(map[string]struct{})

	for _, item := range s {
		if _, ok := uniqueMap[item]; !ok {
			uniqueMap[item] = struct{}{}
		}
	}

	result := make([]string, 0, len(uniqueMap))

	for item := range uniqueMap {
		result = append(result, item)
	}

	return result
}
