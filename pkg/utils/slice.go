package utils

func RemoveIndex(s []int, index int) []int {
	return append(s[:index], s[index+1:]...)
}
