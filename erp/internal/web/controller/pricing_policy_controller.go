package controller

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"go-security/erp/internal/service"
	baseApp "go-security/security/service"
	web "go-security/security/web/middleware"
	"net/http"
)

type PricingPolicyController struct {
	UserService           *baseApp.UserService
	Router                *echo.Group
	PricingPolicyService  *service.PricingPolicyService
	ProfilePricingService *service.ProfilePricingService
}

func NewPricingPolicyController(
	router *echo.Group,
	userService *baseApp.UserService,
	pricingPolicyService *service.PricingPolicyService,
	profilePricingService *service.ProfilePricingService,
) *PricingPolicyController {
	return &PricingPolicyController{
		Router:                router,
		UserService:           userService,
		PricingPolicyService:  pricingPolicyService,
		ProfilePricingService: profilePricingService,
	}
}

func (controller *PricingPolicyController) RegisterRoutes() {
	superAdmin, err := controller.UserService.GetRoleByName(context.TODO(), baseApp.RoleSuperAdmin)
	if err != nil {
		panic(err)
	}
	controller.Router.GET("/private/pricing_policy", web.RoleRequired(superAdmin, controller.GetAllPricingPolicies))
	fmt.Println()
}

func (controller *PricingPolicyController) GetAllPricingPolicies(ctx echo.Context) error {
	policies, err := controller.PricingPolicyService.GetAllPolicies(ctx.Request().Context())
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, policies)
}

func (controller *PricingPolicyController) ApplyPricingPolicyToUserProfile(ctx echo.Context) error {
	var request struct {
		ProfileID uint `json:"profile_id"`
		PolicyID  uint `json:"policy_id"`
	}
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	err := controller.ProfilePricingService.ApplyPricingPolicyToProfile(ctx.Request().Context(), request.ProfileID, request.PolicyID)
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusOK)

}
