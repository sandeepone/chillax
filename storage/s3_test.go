package storage

import (
	"github.com/chillaxio/chillax/libenv"
	"os"
	"testing"
)

func S3StorageForTest(t *testing.T) *S3 {
	chillaxEnv := "test"
	chillaxS3AccessKey := libenv.EnvWithDefault("CHILLAX_S3_ACCESS_KEY", "")
	chillaxS3SecretKey := libenv.EnvWithDefault("CHILLAX_S3_SECRET_KEY", "")
	chillaxS3Region := "us-east-1"
	chillaxS3Bucket := "chillax-test"

	if chillaxS3AccessKey == "" || chillaxS3SecretKey == "" {
		t.Fatal("You must set CHILLAX_S3_ACCESS_KEY & CHILLAX_S3_SECRET_KEY environments to run these tests.")
	}

	return NewS3(chillaxEnv, chillaxS3AccessKey, chillaxS3SecretKey, chillaxS3Region, chillaxS3Bucket)
}

func TestS3RootWithDefaultEnvironment(t *testing.T) {
	os.Setenv("CHILLAX_ENV", "test")
	os.Setenv("CHILLAX_STORAGE_TYPE", "s3")

	storage := NewStorage()

	if storage.GetRoot() != "chillax-test" {
		t.Errorf("Root of S3 storage should be located at chillax-test. storage.GetRoot(): %v", storage.GetRoot())
	}
}

func TestS3CreateGetDelete(t *testing.T) {
	storage := S3StorageForTest(t)

	err := storage.Create("/dostuff", []byte("dostuff"))
	if err != nil {
		t.Errorf("Create should not fail. Error: %v, Bucket url: %v", err, storage.Bucket.URL("/dostuff"))
	}

	data, err := storage.Get("/dostuff")
	if err != nil {
		t.Errorf("Get should not fail. Error: %v", err)
	}
	if string(data) != "dostuff" {
		t.Errorf("Get should not fail.")
	}

	err = storage.Delete("/dostuff")
	if err != nil {
		t.Errorf("Delete should not fail. Error: %v", err)
	}
}
