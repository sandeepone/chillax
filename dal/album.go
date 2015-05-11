package dal

import (
	"errors"
	"fmt"
	"github.com/chillaxio/chillax/storage"
	"time"
)

func NewAlbum(storages *storage.Storages, name string) (*Album, error) {
	var err error

	w := &Album{}
	w.storages = storages
	w.bucketName = "albums"
	w.ID = fmt.Sprintf("%v", time.Now().UnixNano())
	w.Name = name

	return w, err
}

type Album struct {
	BaseKV
	Name string
}

func (w *Album) GetBucketName() string {
	return w.bucketName
}

func (w *Album) GetStorages() *storage.Storages {
	return w.storages
}

func (w *Album) ValidateBeforeSave() error {
	if w.Name == "" {
		return errors.New("Name should not be empty.")
	}
	return nil
}

func (w *Album) Save() error {
	return SaveByKey("ID", w.ID, w)
}

func (w *Album) Delete() error {
	return DeleteByKey("ID", w.ID, w)
}
