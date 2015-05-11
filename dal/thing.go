package dal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chillaxio/chillax/libfile"
	"github.com/chillaxio/chillax/storage"
	"mime/multipart"
	"path"
	"time"
)

func NewThing(storages *storage.Storages, userId string) (*Thing, error) {
	var err error

	t := &Thing{}
	t.storages = storages
	t.CreatedAt = time.Now()
	t.ID = fmt.Sprintf("%v", t.CreatedAt.UnixNano())
	t.UserId = userId
	t.Path = t.generatePath()
	t.Mime = "text/plain"

	_, ok := storages.FileSystems[userId]
	if !ok {
		storages.FileSystems[userId] = storage.NewFileSystem(userId)
	}

	return t, err
}

type Thing struct {
	CreatedAt time.Time
	ID        string
	UserId    string
	Path      string
	Mime      string
	storages  *storage.Storages
}

func (t *Thing) generatePath() string {
	if t.ID == "" {
		return ""
	}

	year := t.CreatedAt.Year()
	month := t.CreatedAt.Month()
	day := t.CreatedAt.Day()

	return path.Join(
		fmt.Sprintf("%v", year),
		fmt.Sprintf("%02d", month),
		fmt.Sprintf("%02d", day),
		t.ID,
	)
}

func (t *Thing) generateMetaJson() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Thing) ValidateBeforeSave() error {
	if t.ID == "" {
		return errors.New("ID should not be empty.")
	}
	return nil
}

func (t *Thing) SaveFile(file multipart.File, fileHeader *multipart.FileHeader) error {
	defer file.Close()

	// contentType := fileHeader.Header.Get("Content-Type")

	return nil
}

func (t *Thing) Save(data []byte) error {
	metaJson, err := t.generateMetaJson()
	if err != nil {
		return err
	}

	err = t.storages.FileSystems[t.UserId].Update(path.Join(t.Path, "meta.json"), metaJson)
	if err != nil {
		return err
	}

	ext := libfile.GetExtensionByMime(t.Mime)

	err = t.storages.FileSystems[t.UserId].Update(path.Join(t.Path, t.ID+ext), data)
	if err != nil {
		return err
	}

	return nil
}

func (t *Thing) GetContent() ([]byte, error) {
	ext := libfile.GetExtensionByMime(t.Mime)
	return t.storages.FileSystems[t.UserId].Get(path.Join(t.Path, t.ID+ext))
}

func (t *Thing) Delete() error {
	if t.Path == "" {
		return errors.New("Path cannot be empty.")
	}
	return t.storages.FileSystems[t.UserId].Delete(t.Path)
}
