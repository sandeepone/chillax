package dal

import (
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	chillax_storage "github.com/chillaxio/chillax/storage"
)

func NewUser(storages *chillax_storage.Storages) *User {
	u := &User{}
	u.storages = storages
	u.bucketName = "users"
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

func (u *User) SaveByName() error {
	err := u.ValidateBeforeSave()
	if err != nil {
		return err
	}

	inJson, err := json.Marshal(u)
	if err != nil {
		return err
	}

	return u.storages.KeyValue.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(u.bucketName))
		if err != nil {
			return err
		}

		return bucket.Put([]byte("Name:"+u.Name), inJson)
	})
}
