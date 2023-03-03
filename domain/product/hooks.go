package product

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BeforeSave 保存商品之前先生成sku
func (p *Product) BeforeSave(tx *gorm.DB) (err error) {
	p.SKU = uuid.New().String()
	return
}
