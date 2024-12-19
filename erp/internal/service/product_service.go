package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
	"go-security/erp/internal"
	. "go-security/erp/internal/repository"
	"os"
	"time"
)

const DefaultCategoryName = "Uncategorized"

type ProductService struct {
	FileUploadService IFileUploadService
	ProductRepository IProductRepository
}

func NewProductService(productRepository IProductRepository, fileUploadService IFileUploadService) *ProductService {
	return &ProductService{
		ProductRepository: productRepository,
		FileUploadService: fileUploadService,
	}
}

func (service *ProductService) PostConstruct() {
	uncategorized := &ProductCategory{
		Name: DefaultCategoryName,
	}

	err := service.AddCategory(context.Background(), uncategorized)
	if errors.Is(err, internal.CategoryAlreadyExists) {
		log.Info().Msgf("Category %s already exists, skipping creation", DefaultCategoryName)
		return
	}
	if err != nil {
		log.Fatal().Msgf("Failed to create default category: %v", err)
	}

}

func (service *ProductService) validateProduct(product *Product) error {
	if len(product.Name) == 0 {
		return internal.ProductNameRequired
	}

	if product.CategoryID == 0 {
		return internal.ProductCategoryRequired
	}

	return nil
}

func (service *ProductService) validateCategory(category *ProductCategory) error {
	if category.Name == "" {
		return internal.CategoryNameRequired
	}

	return nil
}

func (service *ProductService) FindAllProducts(ctx context.Context) ([]*Product, error) {
	products, err := service.ProductRepository.FindAllProducts(ctx)
	if err != nil {
		return nil, err
	}
	for idx := range products {
		product := products[idx]
		url, err := service.FileUploadService.GetFileExpiresIn(ctx, product.ProfilePictureObjectName, time.Minute*1)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to get file URL for product %d", product.ID)
		} else {
			product.ProfilePictureURL = url
		}
	}
	return products, nil

}

func (service *ProductService) AddProduct(
	ctx context.Context,
	productName string,
	description string,
	categoryID uint,
	cost uint,
) error {
	product := &Product{
		Name:        productName,
		Description: description,
		CategoryID:  categoryID,
		Cost:        cost,
	}
	err := service.validateProduct(product)
	if err != nil {
		return err
	}

	return service.ProductRepository.AddProduct(ctx, product)
}

func (service *ProductService) AddCategory(ctx context.Context, category *ProductCategory) error {

	if err := service.validateCategory(category); err != nil {
		return err
	}

	categoryFound, err := service.ProductRepository.FindCategoryByName(ctx, category.Name)
	if err == nil && categoryFound != nil {
		return internal.CategoryAlreadyExists
	}

	return service.ProductRepository.AddCategory(ctx, category)
}

func (service *ProductService) DeleteProduct(ctx context.Context, productID uint) error {
	product, err := service.ProductRepository.FindProductByID(ctx, productID)
	if err != nil {
		return err
	}

	if product.HasProfilePicture() {
		if err := service.FileUploadService.DeleteFile(ctx, product.ProfilePictureObjectName); err != nil {
			log.Warn().Err(err).Msgf("Failed to delete profile picture for product %d", productID)
		}
	}
	return service.ProductRepository.DeleteProduct(ctx, productID)
}

func (service *ProductService) DeleteCategory(ctx context.Context, categoryID uint) error {
	category, err := service.ProductRepository.FindCategoryByID(ctx, categoryID)
	if err != nil {
		return err
	}
	if len(category.Products) > 0 {
		return internal.CategoryContainsProducts
	}

	return service.ProductRepository.DeleteCategory(ctx, categoryID)
}

func (service *ProductService) UpdateProduct(ctx context.Context, productID uint, product *Product) error {
	err := service.validateProduct(product)
	if err != nil {
		return err
	}

	return service.ProductRepository.UpdateProduct(ctx, productID, product)
}

func (service *ProductService) UpdateCategory(ctx context.Context, categoryID uint, category *ProductCategory) error {
	err := service.validateCategory(category)
	if err != nil {
		return err
	}
	return service.ProductRepository.UpdateCategory(ctx, categoryID, category)
}

func (service *ProductService) AssignProductToCategory(ctx context.Context, productID, categoryID uint) error {
	return service.ProductRepository.AssignProductToCategory(ctx, productID, categoryID)
}

func (service *ProductService) FindAllProductCategories(ctx context.Context) ([]*ProductCategory, error) {
	return service.ProductRepository.FindAllProductCategories(ctx)

}

func (service *ProductService) FindProductCategoryByID(ctx context.Context, id uint) (*ProductCategory, error) {
	return service.ProductRepository.FindCategoryByID(ctx, id)
}

func (service *ProductService) FindProductByID(ctx context.Context, id uint) (*Product, error) {
	product, err := service.ProductRepository.FindProductByID(ctx, id)
	if err != nil {
		return nil, err
	}
	product.ProfilePictureURL, _ = service.FileUploadService.GetFileExpiresIn(ctx, product.ProfilePictureObjectName, time.Minute*1)
	return product, nil
}

func (service *ProductService) FindCategoryByName(ctx context.Context, name string) (*ProductCategory, error) {
	return service.ProductRepository.FindCategoryByName(ctx, name)
}

func (service *ProductService) UploadProductProfilePicture(ctx context.Context, productID uint, file *os.File) (*URL, *minio.UploadInfo, error) {
	product, err := service.FindProductByID(ctx, productID)
	if err != nil {
		return nil, nil, err
	}
	if product == nil {
		return nil, nil, internal.ProfileNotCreated
	}

	if product.HasProfilePicture() {
		log.Info().Msgf("Deleting old profile image %s", product.ProfilePictureObjectName)
		if err := service.FileUploadService.DeleteFile(ctx, product.ProfilePictureObjectName); err != nil {
			log.Warn().Err(err).Msg("Failed to delete old profile image")
		}
	}

	objectName := uuid.New().String()

	uploadInfo, err := service.FileUploadService.UploadFile(ctx, objectName, file)
	if err != nil {
		return nil, nil, err
	}
	if err := service.ProductRepository.UpdateProduct(ctx, product.ID, &Product{ProfilePictureObjectName: objectName}); err != nil {
		return nil, nil, err
	}

	stringUrl, err := service.FileUploadService.GetFileExpiresIn(ctx, uploadInfo.Key, 5*time.Minute)
	if err != nil {
		return nil, nil, err
	}
	url := URL(stringUrl)
	return &url, uploadInfo, nil
}

func (service *ProductService) GetProductProfileImage(ctx context.Context, productID uint) (*URL, error) {
	product, err := service.FindProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	url, _ := service.FileUploadService.GetFileExpiresIn(ctx, product.ProfilePictureObjectName, 30*time.Minute)
	return (*URL)(&url), nil
}
