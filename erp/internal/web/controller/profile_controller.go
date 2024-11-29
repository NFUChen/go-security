package controller

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"go-security/erp/internal"
	"go-security/erp/internal/repository"
	"go-security/erp/internal/service"
	baseApp "go-security/security/service"
	baseController "go-security/security/web/controller"
	web "go-security/security/web/middleware"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
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
	superAdmin, err := controller.UserService.GetRoleByName(context.TODO(), baseApp.RoleSuperAdmin)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to get super admin role")
	}
	controller.Router.GET("/private/profile-by-id", web.RoleRequired(superAdmin, controller.GetProfileByUserID))
	controller.Router.GET("/private/personal-profile", controller.GetProfile)
	controller.Router.POST("/private/profile", controller.AddProfile)
	controller.Router.GET("private/profile", controller.GetAllProfiles)
	controller.Router.GET("/private/is_complete_profile", controller.IsCompleteProfile)
	controller.Router.POST("/private/self_upload_profile_picture", controller.SelfUploadProfilePicture)
	controller.Router.POST("/private/admin_upload_profile_picture", web.RoleRequired(superAdmin, controller.AdminUploadProfilePicture))
}

func (controller *ProfileController) IsCompleteProfile(ctx echo.Context) error {
	user, _ := baseController.ExtractUserClaims(ctx)
	isCompleted := controller.ProfileService.IsProfileExists(ctx.Request().Context(), user.ID)

	return ctx.JSON(200, map[string]bool{"is_profile_completed": isCompleted})
}

func (controller *ProfileController) GetAllProfiles(ctx echo.Context) error {
	profiles, err := controller.ProfileService.GetAllProfiles(ctx.Request().Context())
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, profiles)
}

func (controller *ProfileController) GetProfile(ctx echo.Context) error {
	user, _ := baseController.ExtractUserClaims(ctx)
	profile, err := controller.ProfileService.FindProfileByUserId(ctx.Request().Context(), user.ID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, profile)
}

func (controller *ProfileController) AddProfile(ctx echo.Context) error {
	user, _ := baseController.ExtractUserClaims(ctx)
	var profile repository.UserProfile
	if err := ctx.Bind(&profile); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	err := controller.ProfileService.AddProfile(ctx.Request().Context(), user.ID, profile.PhoneNumber)
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusCreated)
}

func (controller *ProfileController) MultiPartFileToOsFile(src *multipart.FileHeader) (*os.File, error) {
	targetFile, err := os.Create(src.Filename)
	if err != nil {
		return nil, err
	}
	sourceFile, err := src.Open()
	if err != nil {
		return nil, err
	}
	defer sourceFile.Close()

	if _, err := io.Copy(targetFile, sourceFile); err != nil {
		return nil, err
	}
	return targetFile, nil

}

func (controller *ProfileController) UploadHandler(ctx echo.Context, userID uint) error {
	fileFromForm, err := ctx.FormFile("profile_picture")
	if err != nil {
		return err
	}

	osFile, err := controller.MultiPartFileToOsFile(fileFromForm)
	if err != nil {
		return internal.UnableToConvertFile
	}

	defer osFile.Close()

	url, uploadInfo, err := controller.ProfileService.UploadUserProfilePicture(ctx.Request().Context(), uint(userID), osFile)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, map[string]any{"upload_info": uploadInfo, "url": url})
}

func (controller *ProfileController) SelfUploadProfilePicture(ctx echo.Context) error {
	user, _ := baseController.ExtractUserClaims(ctx)
	return controller.UploadHandler(ctx, user.ID)
}

func (controller *ProfileController) AdminUploadProfilePicture(ctx echo.Context) error {
	userID, err := strconv.Atoi(ctx.FormValue("user_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "user_id must be a number"})
	}
	return controller.UploadHandler(ctx, uint(userID))
}

func (controller *ProfileController) GetProfileByUserID(ctx echo.Context) error {

	stringUserID := ctx.QueryParam("user_id")
	if stringUserID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "user_id is required"})
	}

	userID, err := strconv.Atoi(stringUserID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "user_id must be a number"})
	}

	profile, err := controller.ProfileService.FindProfileByUserId(ctx.Request().Context(), uint(userID))
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, profile)
}
