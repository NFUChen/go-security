package controller

import (
	"github.com/labstack/echo/v4"
	"go-security/erp/internal/repository"
	"go-security/erp/internal/service"
	baseApp "go-security/security/service"
	baseController "go-security/security/web/controller"
)

type ProfileController struct {
	Router         *echo.Group
	UserService    *baseApp.UserService
	ProfileService *service.ProfileService
}

func NewProfileController(routerGroup *echo.Group, userService *baseApp.UserService, profileService *service.ProfileService) *ProfileController {
	return &ProfileController{
		UserService:    userService,
		Router:         routerGroup,
		ProfileService: profileService,
	}
}

func (controller *ProfileController) RegisterRoutes() {
	controller.Router.GET("/private/profile", controller.GetProfile)
	controller.Router.POST("/private/profile", controller.AddProfile)
	controller.Router.GET("/private/is_complete_profile", controller.IsCompleteProfile)
}

func (controller *ProfileController) IsCompleteProfile(ctx echo.Context) error {
	user, _ := baseController.ExtractUserClaims(ctx)
	isCompleted := controller.ProfileService.IsProfileExists(ctx.Request().Context(), user.ID)

	return ctx.JSON(200, map[string]bool{"is_profile_completed": isCompleted})
}

func (controller *ProfileController) GetProfile(ctx echo.Context) error {
	user, _ := baseController.ExtractUserClaims(ctx)
	profile, err := controller.ProfileService.FindProfileByCustomerId(ctx.Request().Context(), user.ID)
	if err != nil {
		return ctx.JSON(404, map[string]string{"error": "Profile not found"})
	}
	return ctx.JSON(200, profile)
}

func (controller *ProfileController) AddProfile(ctx echo.Context) error {
	user, _ := baseController.ExtractUserClaims(ctx)
	var profile repository.UserProfile
	if err := ctx.Bind(&profile); err != nil {
		return ctx.JSON(400, map[string]string{"error": "Invalid request"})
	}

	err := controller.ProfileService.AddProfile(ctx.Request().Context(), user.ID, profile.NotificationApproaches, profile.PhoneNumber)
	if err != nil {
		return ctx.JSON(500, map[string]string{"error": "Failed to add profile"})
	}
	return ctx.JSON(200, map[string]string{"message": "Profile added"})
}
