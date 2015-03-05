package storage

import (
	"github.com/boltdb/bolt"
	"github.com/chillaxio/chillax/libenv"
	"github.com/chillaxio/chillax/libstring"
	"os"
	"path"
)

func NewDataDir() (string, error) {
	var err error

	path := libenv.EnvWithDefault("CHILLAX_DATA_DIR", "~/chillax")
	path = libstring.ExpandTildeAndEnv(path)

	err = os.MkdirAll(libstring.ExpandTildeAndEnv(path), 0755)
	if err != nil {
		return "", err
	}

	return path, nil
}

func NewStorages() (*Storages, error) {
	dataDir, err := NewDataDir()
	if err != nil {
		return nil, err
	}

	db, err := bolt.Open(path.Join(dataDir, "kv-db"), 0644, nil)
	if err != nil {
		return nil, err
	}

	s := &Storages{}
	s.DataDir = dataDir
	s.KeyValue = db

	return s, nil
}

type Storages struct {
	DataDir  string
	KeyValue *bolt.DB
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
