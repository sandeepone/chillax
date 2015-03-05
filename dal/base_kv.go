package dal

import (
	chillax_storage "github.com/chillaxio/chillax/storage"
)

type BaseKV struct {
	bucketName string
	storages   *chillax_storage.Storages
}
