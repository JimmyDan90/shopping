package category

// CreateCategoryRequest 分类请求参数
type CreateCategoryRequest struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

// CreateCategoryResponse 创建分类响应参数类型
type CreateCategoryResponse struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}
