package cart

import (
	"gorm.io/gorm"
	"shopping/domain/product"
	"shopping/domain/user"
)

type Cart struct {
	gorm.Model
	UserID uint
	User   user.User
}

// NewCart 实例化
func NewCart(uid uint) *Cart {
	return &Cart{
		UserID: uid,
	}
}

type Item struct {
	gorm.Model
	Product   product.Product
	ProductID uint
	Count     int
	CartID    uint
	Cart      Cart
}

func NewCartItem(productId, cartId uint, count int) *Item {
	return &Item{
		ProductID: productId,
		Count:     count,
		CartID:    cartId,
	}
}
