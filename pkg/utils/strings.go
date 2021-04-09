package utils

func UniqueStrings(strs []string) []string {
	var result []string
	var strMap = map[string]bool{}
	for _, str := range strs {
		if strMap[str] {
			continue
		}
		strMap[str] = true
		result = append(result, str)
	}
	return result
}
