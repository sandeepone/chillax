package dal

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	chillax_storage "github.com/chillaxio/chillax/storage"
)

type BaseKV struct {
	ID         string
	bucketName string
	storages   *chillax_storage.Storages
}

func (b *BaseKV) SaveByKey(key, value string, validationFunc func() error) error {
	err := validationFunc()
	if err != nil {
		return err
	}

	inJson, err := json.Marshal(b)
	if err != nil {
		return err
	}

	return b.storages.KeyValue.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(b.bucketName))
		if err != nil {
			return err
		}

		return bucket.Put([]byte(fmt.Sprintf("%v:%v", key, value)), inJson)
	})
}
