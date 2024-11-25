package controller

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	"github.com/rs/zerolog/log"
	. "go-security/erp/internal/service/notification"
)

type LineController struct {
	ChannelSecret string
	Router        *echo.Group
	LineService   *LineService
}

func (controller *LineController) RegisterRoutes() {
	controller.Router.POST("/public/line/callback", controller.LineCallbackHandler)
}

func (controller *LineController) LineCallbackHandler(context echo.Context) error {

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

func NewLineController(router *echo.Group, channelSecret string, lineService *LineService) *LineController {
	return &LineController{
		ChannelSecret: channelSecret,
		Router:        router,
		LineService:   lineService,
	}
}
