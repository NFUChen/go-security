package view

import (
	. "go-security/erp/internal/repository"
)

type HeaderType string

const (
	HeaderTypeText        HeaderType = "text"
	HeaderTypeNumber      HeaderType = "number"
	HeaderTypeLongText    HeaderType = "long_text"
	HeaderTypeImageSource HeaderType = "image_source"
)

type TableHeader struct {
	Key   string     `json:"key"`   // Unique key for the column
	Label string     `json:"label"` // Display name for the column
	Type  HeaderType `json:"type"`  // Type of data in the column
}

type Table[T any] struct {
	Headers []*TableHeader `json:"headers"` // Column headers
	Rows    []*T           `json:"rows"`    // Rows of data
}

type TableService struct{}

func NewTableService() *TableService {
	return &TableService{}
}

func (service *TableService) generateProductCategoryHeaders() []*TableHeader {
	return []*TableHeader{
		{
			Key:   "id",
			Label: "索引",
			Type:  HeaderTypeText,
		},
		{
			Key:   "name",
			Label: "產品類別名稱",
			Type:  HeaderTypeText,
		},
		{
			Key:   "description",
			Label: "產品類別描述",
			Type:  HeaderTypeLongText,
		},
	}
}

func (service *TableService) generateProductCategoryRows(categories []*ProductCategory) []*ProductCategory {
	var rows []*ProductCategory
	for _, category := range categories {
		rows = append(rows, category)
	}
	return rows
}

func (service *TableService) generateProductTableHeaders() []*TableHeader {
	return []*TableHeader{
		{
			Key:   "id",
			Label: "索引",
			Type:  HeaderTypeText,
		},
		{
			Key:   "name",
			Label: "產品名稱",
			Type:  HeaderTypeText,
		},
		{
			Key:   "description",
			Label: "產品描述",
			Type:  HeaderTypeLongText,
		},
		{
			Key:   "profilePictureObjectName",
			Label: "產品圖片",
			Type:  HeaderTypeImageSource,
		},
	}
}

func (service *TableService) generateProductTableRows(products []*Product) []*Product {
	var rows []*Product
	for _, product := range products {
		rows = append(rows, product)
	}
	return rows
}

func (service *TableService) GenerateProductTable(products []*Product) *Table[Product] {
	return &Table[Product]{
		Headers: service.generateProductTableHeaders(),
		Rows:    service.generateProductTableRows(products),
	}
}

func (service *TableService) GenerateProductCategoryTable(categories []*ProductCategory) *Table[ProductCategory] {
	return &Table[ProductCategory]{
		Headers: service.generateProductCategoryHeaders(),
		Rows:    service.generateProductCategoryRows(categories),
	}
}
