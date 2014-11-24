package storage

import (
	"github.com/chillaxio/chillax/libenv"
)

type Storer interface {
	GetRoot() string
	Create(string, []byte) error
	Update(string, []byte) error
	Get(string) ([]byte, error)
	List(string) ([]string, error)
	Delete(string) error
}

func NewStorage() Storer {
	storageType := libenv.EnvWithDefault("STORAGE_TYPE", "FileSystem")
	chillaxEnv := libenv.EnvWithDefault("CHILLAX_ENV", "development")

	if storageType == "FileSystem" {
		return NewFileSystem(chillaxEnv)
	}
	if storageType == "S3" {
		chillaxS3AccessKey := libenv.EnvWithDefault("CHILLAX_S3_ACCESS_KEY", "")
		chillaxS3SecretKey := libenv.EnvWithDefault("CHILLAX_S3_SECRET_KEY", "")
		chillaxS3Region := libenv.EnvWithDefault("CHILLAX_S3_REGION", "")
		chillaxS3Bucket := libenv.EnvWithDefault("CHILLAX_S3_BUCKET", "")

		return NewS3(chillaxEnv, chillaxS3AccessKey, chillaxS3SecretKey, chillaxS3Region, chillaxS3Bucket)
	}
	return nil
}
