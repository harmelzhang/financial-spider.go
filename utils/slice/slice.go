package slice

// 查询指定元素在数组中的位置
func IndexOf[T string](source []T, target T) int {
	for index, item := range source {
		if item == target {
			return index
		}
	}
	return -1
}

// 根据每页数量对数组进行切片
func ArraySlice[T any](source []T, pageSize int) ([][]T, int) {
	result := make([][]T, 0)

	totalPages := len(source) / pageSize
	if len(source)%pageSize != 0 {
		totalPages++
	}

	for i := 1; i <= totalPages; i++ {
		start := (i - 1) * pageSize
		end := start + pageSize
		if end > len(source) {
			end = len(source)
		}
		result = append(result, source[start:end])
	}

	return result, len(result)
}
