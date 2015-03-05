package dal

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	chillax_storage "github.com/chillaxio/chillax/storage"
	"time"
)

type BaseKV struct {
	ID         string
	bucketName string
	storages   *chillax_storage.Storages
}

func (b *BaseKV) ValidateBeforeSave() error {
	return nil
}

func (b *BaseKV) Save() error {
	err := b.ValidateBeforeSave()

	if b.ID == "" {
		b.ID = fmt.Sprintf("%v", time.Now().UnixNano())
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

		return bucket.Put([]byte("ID:"+b.ID), inJson)
	})
}
