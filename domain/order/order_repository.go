package order

import (
	"fmt"
	"gorm.io/gorm"
	"log"
)

type Repository struct {
	db *gorm.DB
}

// NewOrderRepository 实例化
func NewOrderRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Migration 创建表
func (r *Repository) Migration() {
	err := r.db.AutoMigrate(&Order{})
	if err != nil {
		log.Print(err)
	}
}

// FindByOrderID 根据订单id查找
func (r *Repository) FindByOrderID(oid uint) (*Order, error) {
	var currentOrder *Order
	if err := r.db.Where("IsCanceled = ?", false).Where("ID", oid).First(&currentOrder).Error; err != nil {
		return nil, err
	}
	return currentOrder, nil
}

// Update 更新订单
func (r *Repository) Update(newOrder Order) error {
	result := r.db.Save(&newOrder)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Create 创建订单
func (r *Repository) Create(ci *Order) error {
	result := r.db.Create(ci)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetAll 获得所有的订单
func (r *Repository) GetAll(pageIndex, pageSize int, uid uint) ([]EachOrder, int) {
	var orders []Order
	var orderItems []OrderedItem
	eachOrders := make([]EachOrder, 0, 1000)
	var count int64
	r.db.Where("IsCanceled = ?", 0).Where(
		"UserID", uid).Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&orders).Count(&count)
	for _, order := range orders {
		items := r.db.Where("OrderID = ?", order.ID).Find(&orderItems)
		fmt.Println("items: ", items)
		eachOrders = append(eachOrders, EachOrder{
			Order:        order,
			OrderedItems: orderItems,
		})
		//	//r.db.Where("OrderID = ?", order.ID).Find(&orders[i].OrderedItems)
		//	//for j, item := range orders[i].OrderedItems {
		//	//	r.db.Where("ID = ?", item.ProductID).First(&orders[i].OrderedItems[j].Product)
		//	//}
	}
	return eachOrders, int(count)
}
