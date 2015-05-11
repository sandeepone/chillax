package storage

import (
	"github.com/boltdb/bolt"
	"github.com/chillaxio/chillax/libenv"
	"github.com/chillaxio/chillax/libstring"
	"github.com/gorilla/sessions"
	"os"
	"path"
)

func NewDataDir(path string) (string, error) {
	var err error

	if path == "" {
		path = "~/chillax"
	}
	path = libenv.EnvWithDefault("CHILLAX_DATA_DIR", path)
	path = libstring.ExpandTildeAndEnv(path)

	err = os.MkdirAll(libstring.ExpandTildeAndEnv(path), 0755)
	if err != nil {
		return "", err
	}

	return path, nil
}

func NewStoragesGivenDataDir(dataDir string) (*Storages, error) {
	kvdb, err := bolt.Open(path.Join(dataDir, "kv-db"), 0644, nil)
	if err != nil {
		return nil, err
	}

	cookieStore, err := NewCookie(kvdb)
	if err != nil {
		return nil, err
	}

	s := &Storages{}
	s.DataDir = dataDir
	s.KeyValue = kvdb
	s.Cookie = cookieStore

	return s, nil
}

func NewStorages() (*Storages, error) {
	dataDir, err := NewDataDir("")
	if err != nil {
		return nil, err
	}
	return NewStoragesGivenDataDir(dataDir)
}

func NewTestStorages() (*Storages, error) {
	dataDir, err := NewDataDir("~/chillax-test")
	if err != nil {
		return nil, err
	}
	return NewStoragesGivenDataDir(dataDir)
}

type Storages struct {
	DataDir  string
	KeyValue *bolt.DB
	Cookie   *sessions.CookieStore
}

func (s *Storages) CreateKVBucket(name string) error {
	return s.KeyValue.Update(func(tx *bolt.Tx) error {
		tx.CreateBucket([]byte(name))
		return nil
	})
}

func (s *Storages) RemoveAll() error {
	err := s.KeyValue.Close()
	if err != nil {
		return err
	}

	return os.RemoveAll(s.DataDir)
}
