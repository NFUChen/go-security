package application

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/rs/zerolog/log"
	"go-security/erp/internal/controller"
	"go-security/erp/internal/repository"
	appService "go-security/erp/internal/service"
	"go-security/erp/internal/service/notification"
	. "go-security/security/application"
	"go-security/security/service"
	baseController "go-security/security/web/controller"
)

func MustNewErpApplicationContext(appConfig *ErpApplicationConfig, baseApp *Application, baseContext *ApplicationContext) *ApplicationContext {
	var smtpService *service.SmtpService
	for _, _service := range baseContext.Services {
		serviceFound, ok := _service.(*service.SmtpService)
		if ok {
			smtpService = serviceFound
			break
		}
	}
	if smtpService == nil {
		log.Fatal().Msgf("Unable to find SMTP service in base context")
	}

	awsConfig, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(appConfig.Aws.Region),
		config.WithCredentialsProvider(
			credentials.StaticCredentialsProvider{
				Value: aws.Credentials{
					AccessKeyID:     appConfig.Aws.AwsAccessKeyID,
					SecretAccessKey: appConfig.Aws.AwsSecretAccessKey},
			},
		),
	)
	if err != nil {
		log.Fatal().Msgf("failed to load configuration, %v", err)
	}

	snsClient := sns.NewFromConfig(awsConfig)
	snsService := notification.NewAwsSnsService(snsClient)

	services := []service.IService{
		snsService,
	}

	orderRepo := repository.NewOrderRepository(baseApp.SqlEngine)
	profileRepo := repository.NewProfileRepository()
	profileService := appService.NewProfileService(profileRepo)
	emailService := notification.NewEmailService(smtpService)
	lineService := notification.NewLineService()
	_ = appService.NewOrderService(orderRepo, profileService, emailService, snsService, lineService)
	router := baseApp.Engine.Group("/erp-api")
	lineController := controller.NewLineController(router, appConfig.Line.ChannelSecret, lineService)

	controllers := []baseController.Controller{
		lineController,
	}

	return &ApplicationContext{
		Models:      repository.NewCoreModelProvider().ProvideModels(),
		Services:    services,
		Controllers: controllers,
	}
}
