package cart

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"shopping/domain/user"
)

type Repository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Migration() {
	err := r.db.AutoMigrate(&Cart{})
	if err != nil {
		log.Print(err)
	}
}

func (r *Repository) Update(cart Cart) error {
	result := r.db.Save(cart)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FindOrCreateByUserID 根据用户ID查找或创建购物车
func (r *Repository) FindOrCreateByUserID(userId uint) (*Cart, error) {
	var cart *Cart
	err := r.db.Where(Cart{UserID: userId}).Attrs(NewCart(userId)).FirstOrCreate(&cart).Error
	if err != nil {
		return nil, err
	}
	return cart, nil
}

// FindByUserID 根据用户ID创建购物车
func (r *Repository) FindByUserID(userId uint) (*Cart, error) {
	var cart *Cart
	err := r.db.Where(Cart{UserID: userId}).Attrs(NewCart(userId)).First(&cart).Error
	if err != nil {
		return nil, user.ErrUserNotFound
	}
	return cart, nil
}

type ItemRepository struct {
	db *gorm.DB
}

func NewCartItemRepository(db *gorm.DB) *ItemRepository {
	return &ItemRepository{
		db: db,
	}
}

func (i *ItemRepository) Migration() {
	err := i.db.AutoMigrate(&Item{})
	if err != nil {
		log.Print(err)
	}
}

func (i *ItemRepository) Update(item Item) error {
	result := i.db.Save(&item)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (i *ItemRepository) FindByID(pid uint, cid uint) (*Item, error) {
	var item *Item
	err := i.db.Where(&Item{ProductID: pid, CartID: cid}).First(&item).Error
	if err != nil {
		return nil, errors.New("cart item not found")
	}
	return item, nil
}

func (r *ItemRepository) Create(ci *Item) error {
	result := r.db.Create(ci)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *ItemRepository) GetItems(cartId uint) ([]Item, error) {
	var cartItems []Item
	err := r.db.Where(&Item{CartID: cartId}).Find(&cartItems).Error
	if err != nil {
		return nil, err
	}
	for i, item := range cartItems {
		err := r.db.Model(item).Association("Product").Find(&cartItems[i].Product)
		if err != nil {
			return nil, err
		}
	}
	return cartItems, nil
}
