package cart

import "gorm.io/gorm"

// AfterUpdate 如果计数为0，则删除商品
func (item *Item) AfterUpdate(tx *gorm.DB) (err error) {
	if item.Count <= 0 {
		return tx.Unscoped().Delete(&item).Error
	}
	return
}
