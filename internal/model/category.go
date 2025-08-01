package model

// 行业分类
type Category struct {
	// 分类类型（证券会、中证）
	Type string
	// 行业代码
	Code string
	// 名称
	Name string
	// 层级
	Level string
	// 显示顺序
	DisplayOrder int
	// 父行业代码
	ParentCode string
}
