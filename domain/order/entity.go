package order

import (
	"gorm.io/gorm"
	"shopping/domain/product"
	"shopping/domain/user"
)

// Order 结构体
type Order struct {
	gorm.Model
	UserID     uint
	User       user.User
	TotalPrice float32
	IsCanceled bool
}

// OrderedItem 结构体
type OrderedItem struct {
	gorm.Model
	Product    product.Product
	ProductID  uint
	Count      int
	OrderIO    uint
	IsCanceled bool
}

type CurrentOrder struct {
	Order
	OrderedItems []OrderedItem
}

// NewOrder 实例化订单
func NewOrder(uid uint, items []OrderedItem) *Order {
	var totalPrice float32 = 0.0
	for _, item := range items {
		totalPrice = item.Product.Price
	}
	return &Order{
		UserID:     uid,
		TotalPrice: totalPrice,
		IsCanceled: false,
	}
}

// NewOrderedItem 实例化订单项
func NewOrderedItem(count int, pid uint) *OrderedItem {
	return &OrderedItem{
		Count:      count,
		ProductID:  pid,
		IsCanceled: false,
	}
}
