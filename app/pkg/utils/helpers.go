package utils

import (
	"strconv"
	"strings"
)

// Convert Преобразует строку 1,2,3 в числовой массив  [1,2,3]
func ConvertStringToArray(param string) []int {
	var arr []int
	sortsSplit := strings.Split(param, ",")
	for _, elem := range sortsSplit {
		param, err := strconv.Atoi(elem)
		if err != nil {
			return []int{}
		}
		if elem == "" {
			continue
		}
		arr = append(arr, param)
	}
	return arr
}
