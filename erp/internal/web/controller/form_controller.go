package controller

import (
	"context"
	"github.com/labstack/echo/v4"
	"go-security/erp/internal/service/view"
	"go-security/security/service"
	web "go-security/security/web/middleware"
	"net/http"
	"strconv"
)

type FormController struct {
	UserService *service.UserService
	Router      *echo.Group
	FormService *view.FormService
}

func NewFormController(routerGroup *echo.Group, formService *view.FormService, userService *service.UserService) *FormController {
	return &FormController{
		UserService: userService,
		FormService: formService,
		Router:      routerGroup,
	}
}

func (controller *FormController) RegisterRoutes() {
	superAdmin, err := controller.UserService.GetRoleByName(context.TODO(), service.RoleSuperAdmin)
	if err != nil {
		panic(err)
	}
	controller.Router.GET("/private/form/profile_form_template", controller.GetProfileFormTemplate)
	controller.Router.GET("/private/form/profile", web.RoleRequired(superAdmin, controller.GetProfileFormByUserID))
}

func (controller *FormController) GetProfileFormTemplate(ctx echo.Context) error {
	form := controller.FormService.GetUserProfileFormTemplate()
	return ctx.JSON(http.StatusOK, form)
}

func (controller *FormController) GetProfileFormByUserID(ctx echo.Context) error {
	stringUserID := ctx.QueryParam("user_id")
	if stringUserID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "user_id is required"})
	}

	userID, err := strconv.Atoi(stringUserID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "user_id must be a number"})
	}

	form, err := controller.FormService.GetUserProfileForm(ctx.Request().Context(), uint(userID))
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, form)
}
