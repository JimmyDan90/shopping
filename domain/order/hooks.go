package order

import (
	"gorm.io/gorm"
	"shopping/domain/cart"
	"shopping/domain/product"
)

// BeforeCreate 创建订单之前，先查找购物车并删除
func (order *Order) BeforeCreate(tx *gorm.DB) (err error) {
	var currentCart cart.Cart
	if err := tx.Where("UserID = ?", order.UserID).First(&currentCart).Error; err != nil {
		return err
	}
	if err := tx.Where("CartID = ?", currentCart.ID).Unscoped().Delete(&cart.Item{}).Error; err != nil {
		return err
	}
	if err := tx.Unscoped().Delete(&currentCart).Error; err != nil {
		return err
	}
	return nil
}

// BeforeSave 保存之前，更新产品库存
func (o *OrderedItem) BeforeSave(tx *gorm.DB) (err error) {
	var currentProduct product.Product
	var currentOrderItem OrderedItem
	if err := tx.Where("ID = ?", o.ProductID).First(&currentProduct).Error; err != nil {
		return err
	}
	reservedStockCount := 0
	if err := tx.Where("ID = ?", o.ID).First(&currentOrderItem).Error; err != nil {
		reservedStockCount = currentOrderItem.Count
	}
	newStockCount := currentProduct.StockCount + reservedStockCount - o.Count
	if newStockCount < 0 {
		return ErrNotEnoughStock
	}
	if err := tx.Model(&currentProduct).Update("StockCount", newStockCount).Error; err != nil {
		return err
	}
	if o.Count == 0 {
		err := tx.Unscoped().Delete(currentOrderItem).Error
		return err
	}
	return
}

// BeforeUpdate 如果订单被取消，金额将返回库存
func (order *Order) BeforeUpdate(tx *gorm.DB) (err error) {
	if order.IsCanceled {
		var orderedItems []OrderedItem
		if err := tx.Where("OrderID = ?", order.ID).Find(&orderedItems).Error; err != nil {
			return err
		}
		for _, item := range orderedItems {
			var currentProduct product.Product
			if err := tx.Where("ID = ?", item.ProductID).First(&currentProduct).Error; err != nil {
				return err
			}
			newStockCount := currentProduct.StockCount + item.Count
			if err := tx.Model(&currentProduct).Update("StockCount", newStockCount).Error; err != nil {
				return err
			}
			if err := tx.Model(&item).Update("IsCanceled", true).Error; err != nil {
				return err
			}
		}
	}
	return
}
