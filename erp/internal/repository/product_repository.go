package repository

import (
	"context"
	"go-security/erp/internal"
	"gorm.io/gorm"
)

type IProductRepository interface {
	FindAllProducts(ctx context.Context) ([]*Product, error)
	FindCategoryByName(ctx context.Context, name string) (*ProductCategory, error)

	AddProduct(ctx context.Context, product *Product) error
	AddCategory(ctx context.Context, category *ProductCategory) error

	DeleteProduct(ctx context.Context, productID uint) error
	DeleteCategory(ctx context.Context, categoryID uint) error

	UpdateProduct(ctx context.Context, productID uint, product *Product) error
	UpdateCategory(ctx context.Context, categoryID uint, category *ProductCategory) error
	AssignProductToCategory(ctx context.Context, productID uint, categoryID uint) error
	FindAllProductCategories(ctx context.Context) ([]*ProductCategory, error)
	FindCategoryByID(ctx context.Context, categoryID uint) (*ProductCategory, error)
	FindProductByID(ctx context.Context, productID uint) (*Product, error)
}

type ProductRepository struct {
	Engine *gorm.DB
}

func NewProductRepository(engine *gorm.DB) *ProductRepository {
	return &ProductRepository{
		Engine: engine,
	}
}

func (repo *ProductRepository) createProductPreloadQuery(ctx context.Context) *gorm.DB {
	return repo.Engine.WithContext(ctx).Preload("Category")
}

func (repo *ProductRepository) createCategoryPreloadQuery(ctx context.Context) *gorm.DB {
	return repo.Engine.WithContext(ctx).Preload("Products")
}

func (repo *ProductRepository) FindCategoryByID(ctx context.Context, categoryID uint) (*ProductCategory, error) {
	var category ProductCategory
	err := repo.createCategoryPreloadQuery(ctx).Where("id = ?", categoryID).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (repo *ProductRepository) FindProductByID(ctx context.Context, productID uint) (*Product, error) {
	var product Product
	err := repo.createProductPreloadQuery(ctx).Where("id = ?", productID).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (repo *ProductRepository) FindAllProductCategories(ctx context.Context) ([]*ProductCategory, error) {
	var categories []*ProductCategory
	if err := repo.createCategoryPreloadQuery(ctx).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (repo *ProductRepository) FindCategoryByName(ctx context.Context, name string) (*ProductCategory, error) {
	var category ProductCategory
	err := repo.Engine.WithContext(ctx).Where("name = ?", name).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (repo *ProductRepository) AssignProductToCategory(ctx context.Context, productID uint, categoryID uint) error {
	err := repo.Engine.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			var product Product
			err := tx.Find(&product, productID)
			if err != nil {
				return internal.ProductNotFound
			}
			var category ProductCategory
			err = tx.Find(&category, categoryID)
			if err != nil {
				return internal.ProductCategoryNotFound
			}
			product.CategoryID = category.ID
			return tx.Model(&product).Updates(product).Error
		},
	)

	return err
}

func (repo *ProductRepository) AddProduct(ctx context.Context, product *Product) error {
	return repo.Engine.WithContext(ctx).Create(product).Error
}

func (repo *ProductRepository) AddCategory(ctx context.Context, category *ProductCategory) error {
	return repo.Engine.WithContext(ctx).Create(category).Error
}

func (repo *ProductRepository) DeleteProduct(ctx context.Context, productID uint) error {
	return repo.Engine.WithContext(ctx).Delete(&Product{}, productID).Error
}

func (repo *ProductRepository) DeleteCategory(ctx context.Context, categoryID uint) error {
	return repo.Engine.WithContext(ctx).Delete(&ProductCategory{}, categoryID).Error
}

func (repo *ProductRepository) UpdateProduct(ctx context.Context, productID uint, product *Product) error {
	var existingProduct Product
	err := repo.createProductPreloadQuery(ctx).Where("id = ?", productID).First(&existingProduct).Error
	if err != nil {
		return internal.ProductNotFound
	}

	return repo.Engine.Model(&existingProduct).Updates(product).Error
}

func (repo *ProductRepository) UpdateCategory(ctx context.Context, categoryID uint, category *ProductCategory) error {
	return repo.Engine.WithContext(ctx).Model(&ProductCategory{}).Where("id = ?", categoryID).Updates(category).Error

}

func (repo *ProductRepository) FindAllProducts(ctx context.Context) ([]*Product, error) {
	var products []*Product
	if err := repo.Engine.WithContext(ctx).Preload("Category").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
