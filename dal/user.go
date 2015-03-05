package dal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	chillax_storage "github.com/chillaxio/chillax/storage"
	"time"
)

func NewUser(storages *chillax_storage.Storages) *User {
	u := &User{}
	u.storages = storages
	u.bucketName = "users"
	u.ID = fmt.Sprintf("%v", time.Now().UnixNano())

	return u
}

func GetUserById(storages *chillax_storage.Storages, id string) (*User, error) {
	u := NewUser(storages)

	err := storages.KeyValue.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(u.bucketName))

		inJson := bucket.Get([]byte("ID:" + id))
		if inJson == nil {
			return nil
		}

		return json.Unmarshal(inJson, &u)
	})

	return u, err
}

type User struct {
	BaseKV
	Name     string
	Password string
}

func (u *User) ValidateBeforeSave() error {
	if u.Name == "" {
		return errors.New("Name should not be empty.")
	}
	if u.Password == "" {
		return errors.New("Password should not be empty.")
	}
	return nil
}

func (u *User) Save() error {
	return u.SaveByKey("ID", u.ID, u.ValidateBeforeSave)
}

func (u *User) SaveByName() error {
	return u.SaveByKey("Name", u.Name, u.ValidateBeforeSave)
}
