package chillax_storage

import (
    "os"
    "path"
    "errors"
    "io/ioutil"
    "github.com/didip/chillax/libenv"
)

func NewStorage() *Storer {
    storage := libenv.EnvWithDefault("STORAGE_TYPE", "FileSystem")
    if storage == "FileSystem" {
        return &FileSystem{}
    }
}

type Storer interface {
    Create(string, []byte) error
    Update(string, []byte) error
    Get(string)            ([]byte, error)
    Delete(string)         error
}

type FileSystem struct {}

func (fs *FileSystem) CreateOrUpdate(fullpath string, data []byte) error {
    var err error

    basepath := path.Dir(fullpath)
    filename := path.Base(fullpath)

    if _, err = os.Stat(fullpath); os.IsNotExist(err) {
        // Create parent directory
        err = os.MkdirAll(basepath, 0744)
        if err != nil { return err }

        // Create file
        fileHandler, err := os.Create(fullpath)
        if err != nil { return err }
        defer fileHandler.Close()
    }

    err = ioutil.WriteFile(fullpath, string(data), 0744)

    return err
}

func (fs *FileSystem) Create(fullpath string, data []byte) error {
    return fs.CreateOrUpdate(fullpath, data)
}

func (fs *FileSystem) Update(fullpath string, data []byte) error {
    return fs.CreateOrUpdate(fullpath, data)
}

func (fs *FileSystem) Get(fullpath string) ([]byte, error) {
    return ioutil.ReadFile(fullpath)
}

func (fs *FileSystem) Delete(fullpath string) error {
    finfo, err := os.Stat(fullpath)
    if err == nil {
        if finfo.IsDir() {
            err = os.RemoveAll(fullpath)
        }
    }
    return err
}