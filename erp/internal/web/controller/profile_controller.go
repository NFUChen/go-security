package controller

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"go-security/erp/internal"
	"go-security/erp/internal/service"
	"go-security/erp/internal/service/view"
	"go-security/erp/internal/web"
	baseApp "go-security/security/service"
	baseController "go-security/security/web/controller"
	baseWeb "go-security/security/web/middleware"
	"net/http"
	"os"
	"strconv"
)

type ProfileController struct {
	Router              *echo.Group
	UserService         *baseApp.UserService
	ProfileService      *service.ProfileService
	NotificationService *service.NotificationApproachService
	FormAdaptor         *view.FormAdaptor
}

func NewProfileController(
	routerGroup *echo.Group,
	userService *baseApp.UserService,
	profileService *service.ProfileService,
	notificationService *service.NotificationApproachService,
	formAdaptor *view.FormAdaptor,
) *ProfileController {
	return &ProfileController{
		UserService:         userService,
		Router:              routerGroup,
		ProfileService:      profileService,
		NotificationService: notificationService,
		FormAdaptor:         formAdaptor,
	}

}

func (controller *ProfileController) RegisterRoutes() {
	superAdmin, err := controller.UserService.GetRoleByName(context.TODO(), baseApp.RoleSuperAdmin)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to get super admin role")
	}
	controller.Router.GET("/private/profile_by_id", baseWeb.RoleRequired(superAdmin, controller.GetProfileByUserID))
	controller.Router.GET("/private/personal_profile", baseWeb.RoleRequired(superAdmin, controller.GetProfileByUserID))
	controller.Router.PUT("/private/profile", controller.UpdateProfile)
	controller.Router.GET("private/profile", baseWeb.RoleRequired(superAdmin, controller.GetAllProfiles))
	controller.Router.GET("/private/is_self_complete_profile", controller.IsSelfCompleteProfile)
	controller.Router.GET("/private/is_complete_profile", controller.IsUserCompleteProfile)
	controller.Router.POST("/private/self_upload_profile_picture", controller.SelfUploadProfilePicture)
	controller.Router.GET("/private/profile_image", baseWeb.RoleRequired(superAdmin, controller.GetProfileImage))
	controller.Router.POST("/private/admin_upload_profile_picture", baseWeb.RoleRequired(superAdmin, controller.AdminUploadProfilePicture))
	controller.Router.POST("/private/create_default_profile", baseWeb.RoleRequired(superAdmin, controller.CreateDefaultProfile))
}

func (controller *ProfileController) IsUserCompleteProfile(ctx echo.Context) error {
	userID, err := web.GetUserIdFromQueryParam(ctx)
	if err != nil {
		return err
	}
	isCompleted := controller.ProfileService.IsProfileExists(ctx.Request().Context(), userID)
	return ctx.JSON(200, map[string]bool{"is_profile_completed": isCompleted})
}

func (controller *ProfileController) IsSelfCompleteProfile(ctx echo.Context) error {
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

func (controller *ProfileController) UpdateProfile(ctx echo.Context) error {
	userID, err := web.GetUserIdFromQueryParam(ctx)
	if err != nil {
		return err
	}

	profileForm := view.Form{}
	if err := ctx.Bind(&profileForm); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	profile, err := controller.FormAdaptor.FormToUserProfile(userID, &profileForm)
	if err != nil {
		return err
	}

	err = controller.ProfileService.UpdateProfile(ctx.Request().Context(), userID, profile)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, profile)
}

func (controller *ProfileController) UploadHandler(ctx echo.Context, userID uint) error {
	fileFromForm, err := ctx.FormFile("profile_picture")
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

	url, uploadInfo, err := controller.ProfileService.UploadUserProfilePicture(ctx.Request().Context(), userID, osFile)
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

	userID, err := web.GetUserIdFromQueryParam(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	profile, err := controller.ProfileService.GetProfileByUserID(ctx.Request().Context(), userID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, profile)
}

func (controller *ProfileController) EnableNotification(ctx echo.Context) error {
	userID, err := web.GetUserIdFromQueryParam(ctx)
	if err != nil {
		return err
	}

	err = controller.NotificationService.EnableUserNotificationForUser(ctx.Request().Context(), userID)
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusOK)
}

func (controller *ProfileController) IsNotificationEnabled(ctx echo.Context) error {
	userID, err := web.GetUserIdFromQueryParam(ctx)
	if err != nil {
		return err
	}

	isEnabled := controller.NotificationService.IsUserNotificationEnabled(ctx.Request().Context(), userID)
	return ctx.JSON(http.StatusOK, map[string]bool{"is_enabled": isEnabled})
}

func (controller *ProfileController) CreateDefaultProfile(ctx echo.Context) error {
	var request struct {
		UserID uint `json:"user_id"`
	}

	if err := ctx.Bind(&request); err != nil {
		return err
	}
	profile, err := controller.ProfileService.CreateDefaultProfile(ctx.Request().Context(), request.UserID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, profile)
}

func (controller *ProfileController) GetProfileImage(ctx echo.Context) error {
	userID, err := web.GetUserIdFromQueryParam(ctx)
	if err != nil {
		return err
	}
	url, err := controller.ProfileService.GetProfileImage(ctx.Request().Context(), userID)

	if err != nil {
		return ctx.JSON(http.StatusOK, map[string]string{"url": ""})
	}
	return ctx.JSON(http.StatusOK, map[string]*service.URL{"url": url})

}
