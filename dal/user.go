package dal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	chillax_storage "github.com/chillaxio/chillax/storage"
	"golang.org/x/crypto/bcrypt"
	"io"
	"time"
)

func NewUserGivenJson(storages *chillax_storage.Storages, jsonBody io.ReadCloser) (*User, error) {
	var userArgs map[string]interface{}

	err := json.NewDecoder(jsonBody).Decode(&userArgs)
	if err != nil {
		return nil, err
	}

	if _, ok := userArgs["Email"]; !ok {
		return nil, errors.New("Email key does not exist.")
	}
	if _, ok := userArgs["Password"]; !ok {
		return nil, errors.New("Password key does not exist.")
	}

	u, err := NewUser(storages, userArgs["Email"].(string), userArgs["Password"].(string))
	if err != nil {
		return nil, err
	}

	return u, nil
}

func NewUser(storages *chillax_storage.Storages, email, password string) (*User, error) {
	var err error

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

func GetUserByEmailAndPasswordJson(storages *chillax_storage.Storages, jsonBody io.ReadCloser) (*User, error) {
	var userArgs map[string]interface{}

	err := json.NewDecoder(jsonBody).Decode(&userArgs)
	if err != nil {
		return nil, err
	}

	if _, ok := userArgs["Email"]; !ok {
		return nil, errors.New("Email key does not exist.")
	}
	if _, ok := userArgs["Password"]; !ok {
		return nil, errors.New("Password key does not exist.")
	}

	u, err := GetUserByEmailAndPassword(storages, userArgs["Email"].(string), userArgs["Password"].(string))
	if err != nil {
		return nil, err
	}

	return u, nil
}

func GetUserByEmailAndPassword(storages *chillax_storage.Storages, email, password string) (*User, error) {
	u, err := NewUser(storages, "", "")

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
	err := SaveByKey("ID", u.ID, u)
	if err != nil {
		return err
	}
	return u.SaveByEmail()
}

func (u *User) SaveByEmail() error {
	return SaveByKey("Email", u.Email, u)
}
