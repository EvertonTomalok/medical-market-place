package helpers

import "github.com/thoas/go-funk"

func WorkerHasAllDocumentsRequired(workerDocumentsIds []int64, documentsRequireds []int64) bool {
	for _, docRequired := range documentsRequireds {
		if !funk.Contains(workerDocumentsIds, docRequired) {
			return false
		}
	}
	return true
}

func RemoveDuplicate[T string | int](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
