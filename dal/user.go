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

func NewUser(storages *chillax_storage.Storages, email, password, passwordAgain string) (*User, error) {
	var err error

	if password != passwordAgain {
		return nil, errors.New("Password and PasswordAgain fields does not match.")
	}

	u := &User{}
	u.storages = storages
	u.bucketName = "users"
	u.ID = fmt.Sprintf("%v", time.Now().UnixNano())
	u.Email = email

	if password != "" {
		u.Password, err = u.HashedPassword(password)
		if err != nil {
			return nil, err
		}
	}

	return u, err
}

func GetUserById(storages *chillax_storage.Storages, id string) (*User, error) {
	u, err := NewUser(storages, "", "", "")

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

func GetUserByEmailAndPassword(storages *chillax_storage.Storages, email, password string) (*User, error) {
	u, err := NewUser(storages, "", "", "")

	err = storages.KeyValue.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(u.bucketName))
		if err != nil {
			return err
		}

		inJson := bucket.Get([]byte("Email:" + email))
		if inJson == nil {
			return nil
		}

		err = json.Unmarshal(inJson, u)
		if err != nil {
			return err
		}

		return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	})

	if u.Email == "" || u.Password == "" {
		return nil, errors.New("Failed to get user.")
	}

	return u, err
}

type User struct {
	BaseKV
	Email    string
	Name     string
	Password string
	Albums   map[string]*Album
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

func (u *User) UpdateCreds(email, password, passwordAgain string) error {
	var err error

	if password != passwordAgain {
		return errors.New("Password and PasswordAgain fields does not match.")
	}

	u.Email = email
	if password != "" {
		u.Password, err = u.HashedPassword(password)
		if err != nil {
			return err
		}
	}

	return u.Save()
}

func (u *User) CreateAlbum(name string) error {
	if u.Albums == nil {
		u.Albums = make(map[string]*Album)
	}

	album, err := NewAlbum(u.storages, name)
	if err != nil {
		return err
	}

	err = album.Save()
	if err != nil {
		return err
	}

	u.Albums[album.ID] = album

	return u.Save()
}

func (u *User) GetAlbumByName(name string) *Album {
	for _, album := range u.Albums {
		if album.Name == name {
			return album
		}
	}

	return nil
}

func (u *User) DeleteAlbum(name string) error {
	album := u.GetAlbumByName(name)
	if album != nil {
		delete(u.Albums, album.ID)

		err := u.Save()
		if err != nil {
			return err
		}

		err = album.Delete()
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *User) ValidateBeforeSave() error {
	if u.Email == "" {
		return errors.New("Email should not be empty.")
	}
	if u.Password == "" {
		return errors.New("Password should not be empty.")
	}
	return nil
}

func (u *User) Save() error {
	if u.ID == "" {
		return errors.New("User ID cannot be empty when saving.")
	}

	err := SaveByKey("ID", u.ID, u)
	if err != nil {
		return err
	}
	return u.SaveByEmail()
}

func (u *User) SaveByEmail() error {
	if u.Email == "" {
		return errors.New("User Email cannot be empty when saving.")
	}
	return SaveByKey("Email", u.Email, u)
}
