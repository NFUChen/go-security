package controller

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"go-security/erp/internal"
	"go-security/erp/internal/service"
	"go-security/erp/internal/service/view"
	web "go-security/erp/internal/web"
	baseApp "go-security/security/service"
	baseWeb "go-security/security/web/middleware"
	"net/http"
	"os"
	"strconv"
)

type ProductController struct {
	Router         *echo.Group
	UserService    *baseApp.UserService
	FormAdaptor    *view.FormAdaptor
	ProductService *service.ProductService
}

func NewProductController(
	routerGroup *echo.Group,
	userService *baseApp.UserService,
	formAdaptor *view.FormAdaptor,
	productService *service.ProductService,
) *ProductController {
	return &ProductController{
		Router:         routerGroup,
		UserService:    userService,
		FormAdaptor:    formAdaptor,
		ProductService: productService,
	}
}

func (controller *ProductController) RegisterRoutes() {
	admin, _ := controller.UserService.GetRoleByName(context.TODO(), baseApp.RoleAdmin)
	controller.Router.GET("/private/product",
		baseWeb.RoleRequired(admin, controller.GetAllProducts),
	)
	controller.Router.POST("/private/product",
		baseWeb.RoleRequired(admin, controller.AddProduct),
	)
	controller.Router.POST("/private/category",
		baseWeb.RoleRequired(admin, controller.AddCategory),
	)
	controller.Router.DELETE("/private/product/:id",
		baseWeb.RoleRequired(admin, controller.DeleteProduct),
	)
	controller.Router.DELETE("/private/category/:id",
		baseWeb.RoleRequired(admin, controller.DeleteCategory),
	)
	controller.Router.PUT("/private/product/:id",
		baseWeb.RoleRequired(admin, controller.UpdateProduct),
	)
	controller.Router.PUT("/private/category/:id",
		baseWeb.RoleRequired(admin, controller.UpdateCategory),
	)
	controller.Router.POST("/private/assign_product_to_category",
		baseWeb.RoleRequired(admin, controller.AssignProductToCategory),
	)

	controller.Router.POST("/private/admin_upload_product_profile_picture",
		baseWeb.RoleRequired(admin, controller.AdminUploadProductProfilePicture),
	)

	controller.Router.GET("/private/product_profile_picture", baseWeb.RoleRequired(admin, controller.GetProductProfileImage))

}

func (controller *ProductController) convertIdToInteger(id string) (uint, error) {
	productID, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("id must be a number")
	}
	return uint(productID), nil
}

func (controller *ProductController) GetAllProducts(ctx echo.Context) error {
	products, err := controller.ProductService.FindAllProducts(ctx.Request().Context())
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, products)
}

func (controller *ProductController) AddProduct(ctx echo.Context) error {
	var categoryForm view.Form
	if err := ctx.Bind(&categoryForm); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	product, err := controller.FormAdaptor.FormToProduct(ctx.Request().Context(), &categoryForm)
	if err != nil {
		return err
	}

	if err := controller.ProductService.AddProduct(ctx.Request().Context(), product.Name, product.Description, product.CategoryID, product.Cost); err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, product)
}

func (controller *ProductController) AddCategory(ctx echo.Context) error {
	var categoryForm view.Form
	if err := ctx.Bind(&categoryForm); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	category, err := controller.FormAdaptor.FormToProductCategory(&categoryForm)
	if err != nil {
		return err
	}
	if err := controller.ProductService.AddCategory(ctx.Request().Context(), category); err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, category)
}

func (controller *ProductController) DeleteProduct(ctx echo.Context) error {
	productID := ctx.Param("id")
	id, err := controller.convertIdToInteger(productID)
	if err != nil {
		return err
	}

	if err := controller.ProductService.DeleteProduct(ctx.Request().Context(), id); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

func (controller *ProductController) DeleteCategory(ctx echo.Context) error {
	categoryID := ctx.Param("id")
	id, err := controller.convertIdToInteger(categoryID)
	if err != nil {
		return err
	}

	if err := controller.ProductService.DeleteCategory(ctx.Request().Context(), id); err != nil {
		return err
	}
	return ctx.NoContent(http.StatusOK)
}

func (controller *ProductController) UpdateProduct(ctx echo.Context) error {
	productID := ctx.Param("id")
	id, err := controller.convertIdToInteger(productID)
	if err != nil {
		return err
	}

	var productForm view.Form
	if err := ctx.Bind(&productForm); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	product, err := controller.FormAdaptor.FormToProduct(ctx.Request().Context(), &productForm)

	if err := controller.ProductService.UpdateProduct(ctx.Request().Context(), id, product); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, product)
}

func (controller *ProductController) UpdateCategory(ctx echo.Context) error {
	categoryID := ctx.Param("id")
	id, err := controller.convertIdToInteger(categoryID)
	if err != nil {
		return err
	}
	var categoryForm view.Form
	if err := ctx.Bind(&categoryForm); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	category, err := controller.FormAdaptor.FormToProductCategory(&categoryForm)
	if err != nil {
		return err
	}

	if err := controller.ProductService.UpdateCategory(ctx.Request().Context(), id, category); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, category)
}

func (controller *ProductController) AssignProductToCategory(ctx echo.Context) error {
	var request struct {
		ProductID  uint `json:"product_id"`
		CategoryID uint `json:"category_id"`
	}
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := controller.ProductService.AssignProductToCategory(ctx.Request().Context(), request.ProductID, request.CategoryID); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

func (controller *ProductController) AdminUploadProductProfilePicture(ctx echo.Context) error {
	fileFromForm, err := ctx.FormFile("profile_picture")
	if err != nil {
		return err
	}

	formProductID := ctx.FormValue("product_id")
	productID, err := controller.convertIdToInteger(formProductID)
	if err != nil {
		return err
	}

	osFile, err := web.MultiPartFileToOsFile(fileFromForm)
	if err != nil {
		return internal.UnableToConvertFile
	}

	defer func() {
		if err := os.Remove(osFile.Name()); err != nil {
			log.Error().Err(err).Msg("Unable to remove file")
		}
		log.Info().Msg("File removed, closing file")
		osFile.Close()
	}()

	url, uploadInfo, err := controller.ProductService.UploadProductProfilePicture(ctx.Request().Context(), productID, osFile)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, map[string]any{"upload_info": uploadInfo, "url": url})
}

func (controller *ProductController) GetProductProfileImage(ctx echo.Context) error {
	productID, err := web.GetProductIdFromQueryParam(ctx)
	if err != nil {
		return err
	}
	url, _ := controller.ProductService.GetProductProfileImage(ctx.Request().Context(), productID)
	return ctx.JSON(http.StatusOK, map[string]*service.URL{"url": url})

}
