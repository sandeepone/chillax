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
	Walls    map[string]*Wall
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

func (u *User) CreateWall(name string) error {
	if u.Walls == nil {
		u.Walls = make(map[string]*Wall)
	}

	wall, err := NewWall(u.storages, name)
	if err != nil {
		return err
	}

	err = wall.Save()
	if err != nil {
		return err
	}

	u.Walls[wall.ID] = wall

	return u.Save()
}

func (u *User) GetWallByName(name string) *Wall {
	for _, wall := range u.Walls {
		if wall.Name == name {
			return wall
		}
	}

	return nil
}

func (u *User) DeleteWall(name string) error {
	wall := u.GetWallByName(name)
	if wall != nil {
		delete(u.Walls, wall.ID)

		err := u.Save()
		if err != nil {
			return err
		}

		err = wall.Delete()
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
