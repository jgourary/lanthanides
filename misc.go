package main

func isInMap(dict map[string]int, key string) bool {
	if _, ok := dict[key]; ok {
		return true
	} else {
		return false
	}
}
