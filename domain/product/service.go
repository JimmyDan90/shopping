package product

import "shopping/utils/pagination"

type Service struct {
	productRepository Repository
}

// NewService 实例化
func NewService(productRepository Repository) *Service {
	productRepository.Migration()
	return &Service{
		productRepository: productRepository,
	}
}

// GetAll 获得所有的商品分页
func (s *Service) GetAll(page *pagination.Pages) *pagination.Pages {
	products, count := s.productRepository.GetAll(page.Page, page.PageSize)
	page.Items = products
	page.TotalCount = count
	return page
}

// CreateProduct 创建商品
func (s *Service) CreateProduct(name, desc string, count int, price float32, cid uint) error {
	newProduct := NewProduct(name, desc, count, price, cid)
	err := s.productRepository.Create(newProduct)
	return err
}

// DeleteProduct 删除商品
func (s *Service) DeleteProduct(sku string) error {
	err := s.productRepository.Delete(sku)
	return err
}

// UpdateProduct 更新商品
func (s *Service) UpdateProduct(product *Product) error {
	err := s.productRepository.Update(*product)
	return err
}

// SearchProduct 查找商品
func (s *Service) SearchProduct(text string, page *pagination.Pages) *pagination.Pages {
	products, count := s.productRepository.SearchByString(text, page.Page, page.PageSize)
	page.Items = products
	page.TotalCount = count
	return page
}
