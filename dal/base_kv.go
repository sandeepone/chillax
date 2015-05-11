package dal

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/chillaxio/chillax/storage"
)

type BaseKV struct {
	ID         string
	bucketName string
	storages   *storage.Storages
}

type IKV interface {
	ValidateBeforeSave() error
	GetBucketName() string
	GetStorages() *storage.Storages
}

func SaveByKey(key, value string, kvThing IKV) error {
	err := kvThing.ValidateBeforeSave()
	if err != nil {
		return err
	}

	inJson, err := json.Marshal(kvThing)
	if err != nil {
		return err
	}

	return kvThing.GetStorages().KeyValue.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(kvThing.GetBucketName()))
		if err != nil {
			return err
		}

		return bucket.Put([]byte(fmt.Sprintf("%v:%v", key, value)), inJson)
	})
}

func DeleteByKey(key, value string, kvThing IKV) error {
	return kvThing.GetStorages().KeyValue.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(kvThing.GetBucketName()))
		if err != nil {
			return err
		}

		return bucket.Delete([]byte(fmt.Sprintf("%v:%v", key, value)))
	})
}
