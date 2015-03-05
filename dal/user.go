package dal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	chillax_storage "github.com/chillaxio/chillax/storage"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func NewUser(storages *chillax_storage.Storages, name, password string) (*User, error) {
	var err error

	u := &User{}
	u.storages = storages
	u.bucketName = "users"
	u.ID = fmt.Sprintf("%v", time.Now().UnixNano())
	u.Name = name

	if password != "" {
		u.Password, err = u.HashedPassword(password)
		if err != nil {
			return nil, err
		}
	}

	return u, err
}

func GetUserById(storages *chillax_storage.Storages, id string) (*User, error) {
	u, err := NewUser(storages, "", "")

	err = storages.KeyValue.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(u.bucketName))

		inJson := bucket.Get([]byte("ID:" + id))
		if inJson == nil {
			return nil
		}

		return json.Unmarshal(inJson, u)
	})

	return u, err
}

func GetUserByNameAndPassword(storages *chillax_storage.Storages, name, password string) (*User, error) {
	u, err := NewUser(storages, "", "")

	err = storages.KeyValue.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(u.bucketName))

		inJson := bucket.Get([]byte("Name:" + name))
		if inJson == nil {
			return nil
		}

		err := json.Unmarshal(inJson, u)
		if err != nil {
			return err
		}

		return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	})

	return u, err
}

type User struct {
	BaseKV
	Name     string
	Password string
}

func (u *User) GetBucketName() string {
	return u.bucketName
}

func (u *User) GetStorages() *chillax_storage.Storages {
	return u.storages
}

func (u *User) HashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
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
	err := SaveByKey("ID", u.ID, u)
	if err != nil {
		return err
	}
	return u.SaveByName()
}

func (u *User) SaveByName() error {
	return SaveByKey("Name", u.Name, u)
}
