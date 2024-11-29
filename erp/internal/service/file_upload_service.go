package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

const (
	BucketAlreadyExistsErrorCode = "BucketAlreadyOwnedByYou"
)

type MinioConfig struct {
	Endpoint          string `yaml:"endpoint"`
	AccessKeyID       string `yaml:"access_key_id"`
	SecretAccessKey   string `yaml:"secret_access_key"`
	UseSSL            bool   `yaml:"use_ssl"`
	DefaultBucketName string `yaml:"default_bucket_name"`
}

type IFileUploadService interface {
	UploadFile(ctx context.Context, file *os.File) (*minio.UploadInfo, error)
	DeleteFile(ctx context.Context, objectName string) error
	GetFileExpiresIn(ctx context.Context, objectName string, expiresIn time.Duration) (string, error)
}

type FileUploadService struct {
	Client     *minio.Client
	BucketName string
}

func (service *FileUploadService) PostConstruct() {
	log.Info().Msgf("Creating %s bucket", service.BucketName)
	err := service.Client.MakeBucket(context.Background(), service.BucketName, minio.MakeBucketOptions{})

	var response minio.ErrorResponse
	ok := errors.As(err, &response)
	if ok && response.Code == BucketAlreadyExistsErrorCode {
		log.Info().Msgf("Bucket %s already exists", service.BucketName)
		return
	}
	if err != nil {
		log.Fatal().Msgf("Unable to create %s bucket: %v", service.BucketName, err)
	}
	log.Info().Msgf("Successfully created %s bucket", service.BucketName)
}

func NewFileUploadService(client *minio.Client, bucketName string) *FileUploadService {
	return &FileUploadService{Client: client, BucketName: bucketName}
}

func (service *FileUploadService) DeleteFile(ctx context.Context, objectName string) error {
	err := service.Client.RemoveObject(ctx, service.BucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (service *FileUploadService) UploadFile(ctx context.Context, file *os.File) (*minio.UploadInfo, error) {
	objectName := uuid.New().String()
	ReadUtilROL := int64(-1)
	uploadInfo, err := service.Client.PutObject(ctx, service.BucketName, objectName, file, ReadUtilROL, minio.PutObjectOptions{})
	if err != nil {
		return nil, err
	}
	return &uploadInfo, nil
}

func (service *FileUploadService) GetFileExpiresIn(ctx context.Context, objectName string, expiresIn time.Duration) (string, error) {
	downloadURL, err := service.Client.PresignedGetObject(ctx, service.BucketName, objectName, expiresIn, nil)
	if err != nil {
		return "", err
	}
	return downloadURL.String(), nil
}
