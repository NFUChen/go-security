package controller

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	"github.com/rs/zerolog/log"
	. "go-security/erp/internal/service"
	"go-security/erp/internal/service/notification"
	baseApp "go-security/security/service"
	baseController "go-security/security/web/controller"
	"net/http"
	"time"
)

type LineController struct {
	SecurityConfig *baseApp.SecurityConfig
	ChannelSecret  string

	AuthService      *baseApp.AuthService
	Router           *echo.Group
	LineLoginService *LineLoginService
	LineService      *notification.LineService
}

func (controller *LineController) RegisterRoutes() {
	controller.Router.POST("/public/line/callback", controller.LineCallback)
	controller.Router.POST("/public/line/login", controller.LineLogin)
	controller.Router.GET("/public/line/login-url", controller.GetLineLoginURL)
}

func (controller *LineController) GetLineLoginURL(ctx echo.Context) error {
	url := controller.LineLoginService.GetLineLoginUrl()
	return ctx.JSON(http.StatusOK, map[string]string{"url": url})
}

func (controller *LineController) LineLogin(ctx echo.Context) error {
	var lineAuthRequest struct {
		Code  string `json:"code"`
		State string `json:"state"`
	}

	if err := ctx.Bind(&lineAuthRequest); err != nil {
		log.Warn().Err(err).Msg("Failed to bind request")
		return ctx.NoContent(http.StatusBadRequest)
	}

	authResponse, err := controller.LineLoginService.GetLineAuthResponse(lineAuthRequest.Code)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get line auth response")
		return ctx.NoContent(http.StatusInternalServerError)
	}

	user, err := controller.LineLoginService.VerifyIDToken(authResponse.IDToken)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to verify ID token")
		return ctx.NoContent(http.StatusInternalServerError)
	}

	token, err := controller.LineLoginService.RegisterAndLogin(ctx.Request().Context(), user)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to register and login")
		return ctx.NoContent(http.StatusInternalServerError)
	}
	expiration := time.Until(time.Unix(int64(user.ExpirationTime), 0))
	baseController.WriteCookie(&ctx, baseApp.CookieName, token, expiration)
	return ctx.NoContent(http.StatusOK)

}

func (controller *LineController) LineCallback(context echo.Context) error {

	cb, err := webhook.ParseRequest(controller.ChannelSecret, context.Request())
	if err != nil {
		log.Warn().Err(err).Msg("Failed to parse request")
		if errors.Is(err, webhook.ErrInvalidSignature) {
			return context.NoContent(400)
		} else {
			return context.NoContent(500)
		}
	}

	if err := controller.LineService.ReceiveRequest(cb); err != nil {
		log.Warn().Err(err).Msg("Failed to receive request")
		return context.NoContent(500)
	}
	return nil
}

func NewLineController(router *echo.Group, authService *baseApp.AuthService, lineLoginService *LineLoginService, lineService *notification.LineService, securityConfig *baseApp.SecurityConfig, channelSecret string) *LineController {
	return &LineController{
		Router:           router,
		AuthService:      authService,
		LineLoginService: lineLoginService,
		LineService:      lineService,
		SecurityConfig:   securityConfig,
		ChannelSecret:    channelSecret,
	}
}
