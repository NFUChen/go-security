package controller

import (
	"github.com/labstack/echo/v4"
	"go-security/erp/internal/service"
	"go-security/erp/internal/service/view"
	"net/http"
)

type TableViewController struct {
	Engine         *echo.Group
	TableService   *view.TableService
	ProductService *service.ProductService
}

func (controller *TableViewController) RegisterRoutes() {
	controller.Engine.GET("/private/table/product", controller.GetProductsAsTable)
	controller.Engine.GET("/private/table/product_category", controller.GetProductCategoriesAsTable)
}

func (controller *TableViewController) GetProductsAsTable(ctx echo.Context) error {
	products, err := controller.ProductService.FindAllProducts(ctx.Request().Context())
	if err != nil {
		return err
	}
	table := controller.TableService.GenerateProductTable(products)
	return ctx.JSON(http.StatusOK, table)
}

func (controller *TableViewController) GetProductCategoriesAsTable(ctx echo.Context) error {
	categories, err := controller.ProductService.FindAllProductCategories(ctx.Request().Context())
	if err != nil {
		return err
	}
	table := controller.TableService.GenerateProductCategoryTable(categories)
	return ctx.JSON(http.StatusOK, table)
}

func NewTableViewController(
	router *echo.Group,
	tableService *view.TableService,
	productService *service.ProductService,
) *TableViewController {
	return &TableViewController{
		Engine:         router,
		TableService:   tableService,
		ProductService: productService,
	}
}
