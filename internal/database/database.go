package database

import (
	"internal/config"
	v "internal/version"

	"database/sql"
	"errors"
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
	db  *sql.DB
	cfg *config.Config
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

type ManifestTag struct {
	ID          int64
	uuid        string
	Namespace   string
	Name        string
	Tag         string
	ContentType string
	Sha256      string
	Sha512      string
}

func (mt *ManifestTag) Normalize() error {
	return nil
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	db, err := sql.Open("sqlite3", cfg.DatabaseFile)

	if err != nil {
		return nil, err
	}

	return &Database{
		db:  db,
		cfg: cfg,
	}, nil
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

	_, err = r.db.Exec(query, v.DatabaseVersion())
	if err != nil {
		return err
	}

	query = `
	CREATE TABLE IF NOT EXISTS manifest_tags (
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

func (r *Database) CreateManifestLink(mt ManifestTag) (*ManifestTag, error) {
	query := `
		INSERT
		  INTO manifest_tags
			   (namespace, name, tag, tag_sortable, context_type, sha256, sha512)
		VALUES (        ?,    ?,   ?,            ?,            ?,      ?,      ?)
	`

	sortable := strings.ToLower(mt.Tag)
	res, err := r.db.Exec(query, mt.Namespace, mt.Name, mt.Tag, sortable, mt.ContentType, mt.Sha256, mt.Sha512)
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
	mt.ID = id

	return &mt, nil
}

func (r *Database) GetTags(name string, n int, last string) ([]string, error) {
	query := `
		SELECT tag
		  FROM manifest_tags
		 WHERE namespace = NULL
		   AND name = ?
		`

	if last != "" {
		query += `   AND sortable_tag > ?
			`
	}

	query += " ORDER BY sortable_tag "

	var rows *sql.Rows
	var err error
	if n != 0 {
		query += " LIMIT ?"
		if last != "" {
			rows, err = r.db.Query(query, name, last, n)
		} else {
			rows, err = r.db.Query(query, name, n)
		}
	} else {
		if last != "" {
			rows, err = r.db.Query(query, name, last)
		} else {
			rows, err = r.db.Query(query, name)
		}
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.getColumn(rows, n)
}

func (r *Database) getColumn(rows *sql.Rows, n int) ([]string, error) {
	var answer []string
	if n != 0 {
		answer = make([]string, n)
	} else {
		answer = make([]string, 16)
	}

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		answer = append(answer, name)
	}

	return answer, nil
}

func (r *Database) GetManifestTag(name string, tag string) (*ManifestTag, error) {
	row := r.db.QueryRow(`
		SELECT *
		  FROM manifest_tag
		 WHERE namespace = NULL
		   AND name = ?
		   AND tag  = ?
	`, name, tag)

	var mt ManifestTag
	var st string
	if err := row.Scan(&mt.ID, &mt.uuid, &mt.Namespace, &mt.Name, &mt.Tag, &st, &mt.ContentType, &mt.Sha256, &mt.Sha512); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &mt, nil
}

func (r *Database) GetNamespacedManifestTag(namespace string, name string, tag string) (*ManifestTag, error) {
	row := r.db.QueryRow(`
		SELECT *
		  FROM manifest_tag
		 WHERE namespace = ?
		   AND name = ?
		   AND tag  = ?
	`, namespace, name, tag)

	var mt ManifestTag
	var st string
	if err := row.Scan(&mt.ID, &mt.uuid, &mt.Namespace, &mt.Name, &mt.Tag, &st, &mt.ContentType, &mt.Sha256, &mt.Sha512); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &mt, nil
}


func (r *Database) UpdateManifestTag(id int64, mt ManifestTag) (*ManifestTag, error) {
	if id == 0 {
		return nil, errors.New("invalid updated ID")
	}

	query := "UPDATE manifest_tag SET content_type = ?, sha256 = ?, sha512 = ? WHERE id = ?"
	res, err := r.db.Exec(query, mt.ContentType, mt.Sha256, mt.Sha512, id)
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

	return &mt, nil
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
