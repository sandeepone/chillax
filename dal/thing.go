package dal

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	sqlx_types "github.com/jmoiron/sqlx/types"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func NewThing(db *sqlx.DB) *Thing {
	thing := &Thing{}
	thing.db = db
	thing.table = "things"
	thing.hasID = true

	return thing
}

// CREATE TABLE things (
//     id BIGSERIAL NOT NULL,
//     mime TEXT NOT NULL,
//     size BIGSERIAL NOT NULL,
//     url TEXT NOT NULL,
//     tags JSONB,
//     metadata JSONB
// );
type ThingRow struct {
	ID       int64               `db:"id"`
	Title    string              `db:"title"`
	Mime     string              `db:"mime"`
	Size     int64               `db:"size"`
	URL      string              `db:"url"`
	Tags     sqlx_types.JsonText `db:"tags"`
	Metadata sqlx_types.JsonText `db:"metadata"`
}

type Thing struct {
	Base
}

func (u *Thing) thingRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*ThingRow, error) {
	thingId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return u.GetById(tx, thingId)
}

// AllThings returns all thing rows.
func (u *Thing) AllThings(tx *sqlx.Tx) ([]*ThingRow, error) {
	things := []*ThingRow{}
	query := fmt.Sprintf("SELECT * FROM %v", u.table)
	err := u.db.Select(&things, query)

	return things, err
}

// GetById returns record by id.
func (u *Thing) GetById(tx *sqlx.Tx, id int64) (*ThingRow, error) {
	thing := &ThingRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=$1", u.table)
	err := u.db.Get(thing, query, id)

	return thing, err
}

// CreateFromURL a new record of thing.
// Remember, URL field is for the future saved URL.
func (u *Thing) CreateFromURL(tx *sqlx.Tx, title, url string, tags []string) (*ThingRow, error) {
	if url == "" {
		return nil, errors.New("URL cannot be blank.")
	}

	jsonTags, err := json.Marshal(tags)
	if err != nil {
		return nil, err
	}

	// Do some scraping.

	data := make(map[string]interface{})
	data["id"] = time.Now().UnixNano()
	data["title"] = title
	data["mime"] = "text/plain"
	data["size"] = 0
	data["url"] = ""
	data["tags"] = jsonTags

	sqlResult, err := u.InsertIntoTable(tx, data)
	if err != nil {
		return nil, err
	}

	return u.thingRowFromSqlResult(tx, sqlResult)
}

// CreateFromString a new record of thing.
// Remember, URL field is for the future saved URL.
func (u *Thing) CreateFromString(tx *sqlx.Tx, title, blurb string, tags []string) (*ThingRow, error) {
	if url == "" {
		return nil, errors.New("URL cannot be blank.")
	}

	jsonTags, err := json.Marshal(tags)
	if err != nil {
		return nil, err
	}

	// Upload the blurb text to S3.

	data := make(map[string]interface{})
	data["id"] = time.Now().UnixNano()
	data["title"] = title
	data["mime"] = "text/plain"
	data["size"] = 0
	data["url"] = ""
	data["tags"] = jsonTags

	sqlResult, err := u.InsertIntoTable(tx, data)
	if err != nil {
		return nil, err
	}

	return u.thingRowFromSqlResult(tx, sqlResult)
}
