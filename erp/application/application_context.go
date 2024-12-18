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
	"go-security/erp/internal/service/view"
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

	orderRepo := repository.NewOrderRepository(baseApp.SqlEngine)
	profileRepo := repository.NewProfileRepository(baseApp.SqlEngine)
	notificationApproachRepo := repository.NewNotificationApproachRepository(baseApp.SqlEngine)
	productRepo := repository.NewProductRepository(baseApp.SqlEngine)

	snsClient := sns.NewFromConfig(awsConfig)
	snsService := notification.NewAwsSnsService(snsClient)
	minioClient := appService.MustNewMinIOCredentials(appConfig.Minio)
	fileUploadService := appService.NewFileUploadService(minioClient, appConfig.Minio.DefaultBucketName)
	notificationApproachService := appService.NewNotificationApproachService(notificationApproachRepo)
	productService := appService.NewProductService(productRepo, fileUploadService)

	tableService := view.NewTableService()

	profileService := appService.NewProfileService(baseApp.SqlEngine, appDeps.UserService, profileRepo, fileUploadService, notificationApproachService)
	log.Warn().Msgf("Please make sure to inject pricing policy service to profile service")

	emailService := notification.NewEmailService(appDeps.SmtpService)
	lineService := notification.NewLineService()
	formService := view.NewFormService(profileService, appDeps.UserService, productService, notificationApproachService)
	pricingPolicyRepo := repository.NewPricingPolicyRepository(baseApp.SqlEngine)
	pricingPolicyService := appService.NewPricingPolicyService(pricingPolicyRepo)
	profilePricingService := appService.NewProfilePricingService(profileService, pricingPolicyService) // cross-domain service, for interact with profile and pricing policy

	log.Info().Msg("Injecting pricing policy service to profile service")
	profileService.InjectPricingPolicyService(pricingPolicyService)
	log.Info().Msgf("Injecting profile service to notification approach service")
	notificationApproachService.InjectProfileService(profileService)

	services := []service.IService{
		snsService,
		fileUploadService,
		pricingPolicyService,
		productService,
	}

	_ = appService.NewOrderService(orderRepo, profileService, emailService, snsService, lineService)
	formAdaptor := view.NewFormAdaptor(productService)
	router := baseApp.Engine.Group("/erp-api")
	profileController := controller.NewProfileController(router, appDeps.UserService, profileService, notificationApproachService, formAdaptor)
	lineLoginService := appService.NewLineLoginService(appDeps.AuthService, appDeps.UserService, appConfig.Line)
	lineController := controller.NewLineController(
		router,
		appDeps.AuthService,
		lineLoginService,
		lineService,
		baseApp.AppConfig.Security,
		appConfig.Line.ChannelSecret,
	)
	formController := controller.NewFormController(router, formService, appDeps.UserService, notificationApproachService)

	pricingPolicyController := controller.NewPricingPolicyController(router, appDeps.UserService, pricingPolicyService, profilePricingService)
	productController := controller.NewProductController(router, appDeps.UserService, formAdaptor, productService)

	tableViewController := controller.NewTableViewController(router, tableService, productService)

	controllers := []baseController.Controller{
		lineController,
		profileController,
		formController,
		pricingPolicyController,
		productController,
		tableViewController,
	}

	return &ApplicationContext{
		Models:      repository.NewCoreModelProvider().ProvideModels(),
		Services:    services,
		Controllers: controllers,
	}
}
