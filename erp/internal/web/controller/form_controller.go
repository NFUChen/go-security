package controller

import (
	"context"
	"github.com/labstack/echo/v4"
	service "go-security/erp/internal/service"
	"go-security/erp/internal/service/view"
	web "go-security/erp/internal/web"
	baseApp "go-security/security/service"
	baseWeb "go-security/security/web/middleware"
	"net/http"
)

type ProfileMetaForm struct {
	Form *view.Form     `json:"form"`
	Meta map[string]any `json:"meta"`
}

type FormController struct {
	UserService         *baseApp.UserService
	NotificationService *service.NotificationApproachService
	Router              *echo.Group
	FormService         *view.FormService
}

func NewFormController(routerGroup *echo.Group, formService *view.FormService, userService *baseApp.UserService, notificationService *service.NotificationApproachService) *FormController {
	return &FormController{
		UserService:         userService,
		FormService:         formService,
		Router:              routerGroup,
		NotificationService: notificationService,
	}
}

func (controller *FormController) RegisterRoutes() {
	superAdmin, err := controller.UserService.GetRoleByName(context.TODO(), baseApp.RoleSuperAdmin)
	if err != nil {
		panic(err)
	}
	controller.Router.GET("/private/form/profile_form_template", controller.GetProfileFormTemplate)
	controller.Router.GET("/private/form/product_form_template", controller.GetProductFormTemplate)
	controller.Router.GET("/private/form/product_category_form_template", controller.GetProductCategoryFormTemplate)
	controller.Router.GET("/private/form/product", baseWeb.RoleRequired(superAdmin, controller.GetProductFormByProductID))
	controller.Router.GET("/private/form/profile", baseWeb.RoleRequired(superAdmin, controller.GetProfileFormByUserID))
	controller.Router.GET("/private/form/product_category", baseWeb.RoleRequired(superAdmin, controller.GetProductCategoryFormByCategoryID))
}

func (controller *FormController) GetProfileFormTemplate(ctx echo.Context) error {
	form := controller.FormService.GetUserProfileFormTemplate()
	return ctx.JSON(http.StatusOK, form)
}

func (controller *FormController) GetProfileFormByUserID(ctx echo.Context) error {

	userID, err := web.GetUserIdFromQueryParam(ctx)
	if err != nil {
		return err
	}

	form, err := controller.FormService.GetUserProfileForm(ctx.Request().Context(), userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, form)
}

func (controller *FormController) GetProductFormTemplate(ctx echo.Context) error {
	form, err := controller.FormService.GetProductFormTemplate(ctx.Request().Context())
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, form)
}

func (controller *FormController) GetProductCategoryFormTemplate(ctx echo.Context) error {
	form := controller.FormService.GetProductCategoryFomTemplate()

	return ctx.JSON(http.StatusOK, form)

}

func (controller *FormController) GetProductCategoryFormByCategoryID(ctx echo.Context) error {
	categoryID, err := web.GetCategoryIdFromQueryParam(ctx)
	if err != nil {
		return err
	}

	form, err := controller.FormService.GetProductCategoryForm(ctx.Request().Context(), categoryID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, form)
}

func (controller *FormController) GetProductFormByProductID(ctx echo.Context) error {
	productID, err := web.GetProductIdFromQueryParam(ctx)
	if err != nil {
		return err
	}

	form, err := controller.FormService.GetProductForm(ctx.Request().Context(), productID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, form)

}
