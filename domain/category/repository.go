package category

import (
	"gorm.io/gorm"
	"log"
)

type Repository struct {
	db *gorm.DB
}

// NewCategoryRepository 创建商品分类
func NewCategoryRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Migration 生成商品分类表
func (r *Repository) Migration() {
	err := r.db.AutoMigrate(&Category{})
	if err != nil {
		log.Print(err)
	}
}

// InsertSampleData 生成商品分类测试数据
func (r *Repository) InsertSampleData() {
	categories := []Category{
		{Name: "CAT1", Desc: "Category 1"},
		{Name: "CAT2", Desc: "Category 2"},
	}
	for _, c := range categories {
		r.db.Where(Category{Name: c.Name}).Attrs(Category{
			Name: c.Name,
		}).FirstOrCreate(&c)
	}
}

func (r *Repository) Create(c *Category) error {
	result := r.db.Create(c)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *Repository) GetByName(name string) []Category {
	var categories []Category
	r.db.Where("Name = ?", name).Find(&categories)
	return categories
}

func (r *Repository) BulkCreate(categories []*Category) (int, error) {
	var count int64
	err := r.db.Where(&categories).Count(&count).Error
	return int(count), err
}

func (r *Repository) GetAll(pageIndex, pageSize int) ([]Category, int) {
	var categories []Category
	var count int64
	r.db.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&categories).Count(&count)
	return categories, int(count)
}
