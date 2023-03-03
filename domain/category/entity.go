package category

import "gorm.io/gorm"

// Category 商品分类结构体
type Category struct {
	gorm.Model
	Name     string `gorm:"unique"`
	Desc     string
	IsActive bool
}

// NewCategory 新建商品分类
func NewCategory(name, desc string) *Category {
	return &Category{
		Name:     name,
		Desc:     desc,
		IsActive: true,
	}
}
