package storage

import (
	"fmt"
	goamz_aws "github.com/goamz/goamz/aws"
	goamz_s3 "github.com/goamz/goamz/s3"
	"path"
)

// type Storer interface {
//     GetRoot() string
//     Create(string, []byte) error
//     Update(string, []byte) error
//     Get(string) ([]byte, error)
//     List(string) ([]string, error)
//     Delete(string) error
// }

func NewS3(chillaxEnv, accessKey, secretKey, regionName, bucketName string) *S3 {
	store := &S3{}
	store.ChillaxEnv = chillaxEnv
	store.AccessKey = accessKey
	store.SecretKey = secretKey
	store.RegionName = regionName
	store.BucketName = bucketName
	store.Root = fmt.Sprintf("chillax-%v", chillaxEnv)
	store.Region = goamz_aws.Regions[store.RegionName]

	store.CreateRegion()
	store.CreateConnection()
	store.CreateBucketConnection()
	return store
}

type S3 struct {
	ChillaxEnv string
	AccessKey  string
	SecretKey  string
	RegionName string
	BucketName string
	Root       string
	Region     goamz_aws.Region
	Connection *goamz_s3.S3
	Bucket     *goamz_s3.Bucket
}

func (store *S3) CreateRegion() goamz_aws.Region {
	store.Region = goamz_aws.Regions[store.RegionName]
	store.Region.S3BucketEndpoint = fmt.Sprintf("https://%v.s3.amazonaws.com", store.BucketName)

	return store.Region
}

func (store *S3) CreateConnection() *goamz_s3.S3 {
	auth := goamz_aws.Auth{AccessKey: store.AccessKey, SecretKey: store.SecretKey}
	store.Connection = goamz_s3.New(auth, store.Region)
	return store.Connection
}

func (store *S3) GetRoot() string {
	return store.Root
}

func (store *S3) CreateBucketConnection() *goamz_s3.Bucket {
	store.Bucket = store.Connection.Bucket(store.BucketName)
	return store.Bucket
}

func (store *S3) Create(fullpath string, data []byte) error {
	fullpath = path.Join(store.Root, fullpath)
	return store.Bucket.Put(fullpath, data, "text/plain", goamz_s3.BucketOwnerFull, goamz_s3.Options{})
}

func (store *S3) Update(fullpath string, data []byte) error {
	fullpath = path.Join(store.Root, fullpath)
	return store.Bucket.Put(fullpath, data, "text/plain", goamz_s3.BucketOwnerFull, goamz_s3.Options{})
}

func (store *S3) Get(fullpath string) ([]byte, error) {
	fullpath = path.Join(store.Root, fullpath)
	return store.Bucket.Get(fullpath)
}

func (store *S3) List(fullpath string) ([]string, error) {
	var result []string
	return result, nil
}

func (store *S3) Delete(fullpath string) error {
	fullpath = path.Join(store.Root, fullpath)
	return store.Bucket.Del(fullpath)
}
