package cart

// ItemCartRequest 请求参数
type ItemCartRequest struct {
	SKU   string `json:"sku"`
	Count int    `json:"count"`
}

// CreateCategoryResponse 创建分类响应
type CreateCategoryResponse struct {
	Message string `json:"message"`
}
