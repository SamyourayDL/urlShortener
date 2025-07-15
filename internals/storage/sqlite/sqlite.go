package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"url-shortener/internals/storage"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3" // _ init driver, because we don't use lib explicitly
)

type Storage struct {
	db *sql.DB
}

type StorageErrAlreadyExists struct{}

func (err StorageErrAlreadyExists) Error() string {
	return "Alias already exists! "
}

type StorageErrNoSuchAlias struct{}

func (err StorageErrNoSuchAlias) Error() string {
	return "No such alias! "
}

func New(storagePath string) (*Storage, error) {
	const fn = "storage.sqlite.New" //for error messages, to understand where error occured
	//log.With("in", fn)

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w ", fn, err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS url(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL UNIQUE);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{db}, nil
}

func (s *Storage) SaveURL(urlToSave, alias string) (int64, error) { //(lastInsertedId, err)
	const fn = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url (alias, url) VALUES (?, ?);")
	if err != nil {
		return -1, fmt.Errorf("%s: %w", fn, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(alias, urlToSave)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return -1, fmt.Errorf("%s: %w", fn, storage.ErrURLExists) // url already exists in database
		}

		return -1, fmt.Errorf("%s: %w", fn, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("%s: %w", fn, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const fn = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias=?;")
	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(alias)

	var urlToReturn string
	err = row.Scan(&urlToReturn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", fn, storage.ErrURLNotFound)
		}

		return "", fmt.Errorf("%s: %w", fn, err)
	}

	return urlToReturn, nil
}

func (s *Storage) DeleteURL(alias string) (int64, error) { // (rowsAffected)
	const fn = "storage.sqlite.deleteurl"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias=?;")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(alias)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	rowsDeleted, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return rowsDeleted, nil
}
