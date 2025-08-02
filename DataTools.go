package main

func ConvertToFlattenMatric(matric [][]int8) []int8 {
	result := make([]int8, 0)
	for _, i := range matric {
		for _, j := range i {
			result = append(result, j)
		}
	}
	return result
}
