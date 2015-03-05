package dal

import (
	chillax_storage "github.com/chillaxio/chillax/storage"
)

func NewUser(storages *chillax_storage.Storages) *User {
	u := &User{}
	u.storages = storages
	u.bucketName = "users"
	return u
}

type User struct {
	BaseKV
	Name     string
	Password string
}
