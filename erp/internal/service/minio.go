package service

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
)

func MustNewMinIOCredentials(config *MinioConfig) *minio.Client {
	log.Info().Msgf("Creating MinIO client with endpoint: %s", config.Endpoint)
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})

	if err != nil {
		log.Fatal().Msgf("Unable to create MinIO client: %v", err)
	}
	return client
}
