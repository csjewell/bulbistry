package bulbistry

import (
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/mattn/go-sqlite3"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row not exists")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

type Database struct {
	db *sql.DB
	bc BulbistryConfig
}

type Blob struct {
}

type BlobLink struct {
}

type Manifest struct {
	ID           int64
	uuid         string
	namespace    string
	name         string
	sha          string
	manifestBody string
}

type ManifestLink struct {
	ID          int64
	uuid        string
	Namespace   string
	Name        string
	Tag         string
	ContentType string
	Sha256      string
	Sha512      string
}

func (ml *ManifestLink) Normalize() error {
	return nil
}

func NewDatabase(bc BulbistryConfig) *Database {
	db, err := sql.Open("sqlite3", bc.DatabaseFile)

	if err != nil {
		log.Fatal(err)
	}

	return &Database{
		db: db,
		bc: bc,
	}
}

func (r *Database) InitializeDatabase() error {
	query := `
	CREATE TABLE IF NOT EXISTS tbv_dbversion (
		db_version TEXT
	)`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	query = "INSERT INTO tbv_dbversion (db_version) VALUES (?)"

	_, err = r.db.Exec(query, DatabaseVersion())
	if err != nil {
		return err
	}

	query = `
	CREATE TABLE IF NOT EXISTS manifest_link (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid TEXT NOT NULL UNIQUE,
		namespace TEXT NULL,
		name TEXT NOT NULL,
		tag TEXT NOT NULL,
		tag_sortable TEXT NOT NULL,
		content_type TEXT NOT NULL,
		sha256 TEXT NOT NULL.
		sha512 TEXT NOT NULL
	)`

	_, err = r.db.Exec(query)
	return err
}

func (r *Database) MigrateDatabase() error {
	return nil
}

func (r *Database) CreateManifestLink(ml ManifestLink) (*ManifestLink, error) {
	query := `
		INSERT
		  INTO manifest_link
			   (namespace, name, tag, tag_sortable, context_type, sha256, sha512)
		VALUES (        ?,    ?,   ?,            ?,            ?,      ?,      ?)
	`

	sortable := strings.ToLower(ml.Tag)

	res, err := r.db.Exec(query, ml.Namespace, ml.Name, ml.Tag, sortable, ml.ContentType, ml.Sha256, ml.Sha512)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	ml.ID = id

	return &ml, nil
}

func (r *Database) GetTags(name string, n int, last string) ([]string, error) {
	//row := r.db.QueryRow(`
	//	SELECT *
	//	  FROM websites
	//		 WHERE namespace = NULL
	//		   AND name = ?
	//		   AND tag  = ?
	//	`, name, tag)

}

func (r *Database) GetManifestLink(name string, tag string) (*ManifestLink, error) {
	row := r.db.QueryRow(`
		SELECT *
		  FROM manifest_link
		 WHERE namespace = NULL
		   AND name = ?
		   AND tag  = ?
	`, name, tag)

	var ml ManifestLink
	if err := row.Scan(&ml.ID, &ml.uuid, &ml.Namespace, &ml.Name, &ml.Tag, &ml.ContentType, &ml.Sha256, &ml.Sha512); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &ml, nil
}

func (r *Database) UpdateManifestLink(id int64, updated ManifestLink) (*ManifestLink, error) {
	if id == 0 {
		return nil, errors.New("invalid updated ID")
	}

	res, err := r.db.Exec("UPDATE manifest_link SET content_type = ?, sha256 = ?, sha512 = ? WHERE id = ?", updated.ContentType, updated.Sha256, updated.Sha512, id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, ErrUpdateFailed
	}

	return &updated, nil
}

//func (r *Database) Delete(id int64) error {
//	res, err := r.db.Exec("DELETE FROM websites WHERE id = ?", id)
//	if err != nil {
//		return err
//	}

//	rowsAffected, err := res.RowsAffected()
//	if err != nil {
//		return err
//	}

//	if rowsAffected == 0 {
//		return ErrDeleteFailed
//	}

//	return err
//}
