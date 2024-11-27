package notification

import (
	api "github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	"github.com/rs/zerolog/log"
	. "go-security/erp/internal/repository"
	"net/http"
)

type LineService struct {
	ChannelSecret string
	BlobClient    *api.MessagingApiBlobAPI
	TextClient    *api.MessagingApiAPI
	EventChannel  chan *webhook.MessageEvent
}

func MustNewMessageAPI(channelAccessToken string) *api.MessagingApiAPI {
	client := &http.Client{}
	bot, err := api.NewMessagingApiAPI(
		channelAccessToken,
		api.WithHTTPClient(client),
	)
	if err != nil {
		panic(err)
	}
	return bot
}

func NewLineService() *LineService {
	service := &LineService{

		EventChannel: make(chan *webhook.MessageEvent),
	}

	go service.listenCallbackRequest()

	return service
}

func (service *LineService) ReceiveRequest(request *webhook.CallbackRequest) error {
	for _, event := range request.Events {
		switch eventType := event.(type) {
		case webhook.MessageEvent:
			service.EventChannel <- &eventType
		}
	}
	return nil
}

func (service *LineService) listenCallbackRequest() {
	for {
		message := <-service.EventChannel
		source, _ := message.Source.(webhook.UserSource)
		senderID := source.UserId
		content, _ := message.Message.(webhook.TextMessageContent)
		log.Info().Msgf("Received message: %s from %s", content.Text, senderID)
	}
}

func (service LineService) Name() string {
	return "LineService"
}

func (service LineService) SendOrderWaitingForApprovalMessage(order *CustomerOrder, profile *UserProfile) error {
	//TODO implement me
	panic("implement me")
}

func (service LineService) SendOrderApprovedMessage(order *CustomerOrder, profile *UserProfile) error {
	//TODO implement me
	panic("implement me")
}
