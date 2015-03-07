package dal

import (
	"errors"
	"fmt"
	chillax_storage "github.com/chillaxio/chillax/storage"
	"time"
)

func NewWall(storages *chillax_storage.Storages, name string) (*Wall, error) {
	var err error

	w := &Wall{}
	w.storages = storages
	w.bucketName = "walls"
	w.ID = fmt.Sprintf("%v", time.Now().UnixNano())
	w.Name = name

	return w, err
}

type Wall struct {
	BaseKV
	Name string
}

func (w *Wall) GetBucketName() string {
	return w.bucketName
}

func (w *Wall) GetStorages() *chillax_storage.Storages {
	return w.storages
}

func (w *Wall) ValidateBeforeSave() error {
	if w.Name == "" {
		return errors.New("Name should not be empty.")
	}
	return nil
}

func (w *Wall) Save() error {
	return SaveByKey("ID", w.ID, w)
}

func (w *Wall) Delete() error {
	return DeleteByKey("ID", w.ID, w)
}
