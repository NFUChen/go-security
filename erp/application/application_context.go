package application

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/rs/zerolog/log"
	"go-security/erp/internal/repository"
	appService "go-security/erp/internal/service"
	"go-security/erp/internal/service/notification"
	"go-security/erp/internal/web/controller"
	. "go-security/security/application"
	"go-security/security/service"
	baseController "go-security/security/web/controller"
	"reflect"
)

type dependencies struct {
	SmtpService *service.SmtpService
	AuthService *service.AuthService
	UserService *service.UserService
}

func extractDependencies(context *ApplicationContext) *dependencies {
	var smtpService *service.SmtpService
	var authService *service.AuthService
	var userService *service.UserService
	for _, _service := range context.Services {
		serviceFound, ok := _service.(*service.SmtpService)
		if ok {
			smtpService = serviceFound
		}
		authServiceFound, ok := _service.(*service.AuthService)
		if ok {
			authService = authServiceFound
		}
		userServiceFound, ok := _service.(*service.UserService)
		if ok {
			userService = userServiceFound
		}
	}
	services := []service.IService{
		smtpService, authService, userService,
	}
	for _, _service := range services {
		if _service == nil {
			log.Fatal().Msgf("Unable to find service: %v in base context", reflect.TypeOf(_service))
		}
	}

	return &dependencies{
		SmtpService: smtpService,
		AuthService: authService,
		UserService: userService,
	}
}

func MustNewErpApplicationContext(appConfig *ErpApplicationConfig, baseApp *Application, baseContext *ApplicationContext) *ApplicationContext {
	appDeps := extractDependencies(baseContext)

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
	profileRepo := repository.NewProfileRepository(baseApp.SqlEngine)

	profileService := appService.NewProfileService(profileRepo)
	emailService := notification.NewEmailService(appDeps.SmtpService)
	lineService := notification.NewLineService()

	_ = appService.NewOrderService(orderRepo, profileService, emailService, snsService, lineService)
	router := baseApp.Engine.Group("/erp-api")
	profileController := controller.NewProfileController(router, appDeps.UserService, profileService)
	lineLoginService := appService.NewLineLoginService(appDeps.AuthService, appDeps.UserService, appConfig.Line)
	lineController := controller.NewLineController(
		router,
		appDeps.AuthService,
		lineLoginService,
		lineService,
		baseApp.AppConfig.Security,
		appConfig.Line.ChannelSecret,
	)

	controllers := []baseController.Controller{
		lineController,
		profileController,
	}

	return &ApplicationContext{
		Models:      repository.NewCoreModelProvider().ProvideModels(),
		Services:    services,
		Controllers: controllers,
	}
}
