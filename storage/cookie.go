package storage

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"github.com/boltdb/bolt"
	"github.com/gorilla/sessions"
)

func NewCookie(kvdb *bolt.DB) (*sessions.CookieStore, error) {
	var secret []byte

	// Create cookie secret if one does not exist.
	err := kvdb.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("cookie"))
		if err != nil {
			return err
		}

		secret = bucket.Get([]byte("secret"))
		if len(secret) == 0 {
			secret = []byte(uuid.New())
			err = bucket.Put([]byte("secret"), secret)
		}

		return err
	})

	if err != nil {
		return nil, err
	}

	if len(secret) == 0 {
		return nil, errors.New("Cookie secret cannot be empty.")
	}

	return sessions.NewCookieStore(secret), nil
}
