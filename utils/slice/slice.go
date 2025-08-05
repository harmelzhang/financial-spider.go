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
