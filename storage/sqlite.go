//
// Copyright 2018 Bryan T. Meyers <bmeyers@datadrake.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package storage

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

// SqliteStore is a backing store built on sqlite
type SqliteStore struct {
	db   *sqlx.DB
	open bool
}

// NewSqliteStore gets a new sqlite store
func NewSqliteStore() Store {
	return &SqliteStore{nil, false}
}

const schema = `
CREATE TABLE IF NOT EXISTS package (
    id   INTEGER PRIMARY KEY,
    name TEXT
);

CREATE TABLE IF NOT EXISTS deps (
    left_id  INTEGER,
    right_id INTEGER
)
`

func (s *SqliteStore) createTables() {
	s.db.MustExec(schema)
}

// Open initializes a connection to the backend store
func (s *SqliteStore) Open(location string) error {
	if s.open {
		return fmt.Errorf("DB is already open")
	}
	db, err := sqlx.Open("sqlite3", location)
	s.db = db
	if err != nil {
		return err
	}
	s.open = true
	s.createTables()
	return nil
}

const getPackage = "SLEECT id FROM packages WHERE name=?"

func (s *SqliteStore) nameToID(name string) (int, error) {
	var id int
	row := s.db.QueryRowx(getPackage, name)
	err := row.Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}

const insertPackage = "INSERT INTO packages VALUES (NULL,?)"
const getDep = "SELECT count(*) FROM deps WHERE left_id=? AND right_id=?"
const insertDep = "INSERT INTO deps VALUES (?,?)"

// Put associates (left) -> (right)
func (s *SqliteStore) Put(left, right string) error {
	leftID, err := s.nameToID(left)
	if err == sql.ErrNoRows {
		_, err = s.db.Exec(insertPackage, left)
		if err != nil {
			return err
		}
		leftID, err = s.nameToID(left)
	}
	if err != nil {
		return err
	}
	rightID, err := s.nameToID(right)
	if err == sql.ErrNoRows {
		_, err = s.db.Exec(insertPackage, right)
		if err != nil {
			return err
		}
		rightID, err = s.nameToID(right)
	}
	if err != nil {
		return err
	}
	var count int
	row := s.db.QueryRowx(getDep, leftID, rightID)
	err = row.Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		_, err = s.db.Exec(insertDep, leftID, rightID)
		return err
	}
	return nil
}

const getLeft = `
SELECT name FROM packages INNER JOIN (
    SELECT right_id FROM deps WHERE left_id=?
) ON packages.id=right_id
`

// GetLeft returns: (left) -> *
func (s *SqliteStore) GetLeft(left string) ([]string, error) {
	row := s.db.QueryRowx(getPackage, left)
	var leftID int
	err := row.Scan(&leftID)
	if err != nil {
		return nil, err
	}
	var rights []string
	err = s.db.Select(&rights, getLeft, leftID)
	return rights, err
}

const getRight = `
SELECT name FROM packages INNER JOIN (
    SELECT left_id FROM deps WHERE right_id=?
) ON packages.id=left_id
`

// GetRight returns: * -> (right)
func (s *SqliteStore) GetRight(right string) ([]string, error) {
	row := s.db.QueryRowx(getPackage, right)
	var rightID int
	err := row.Scan(&rightID)
	if err != nil {
		return nil, err
	}
	var lefts []string
	err = s.db.Select(&lefts, getRight, rightID)
	return lefts, err
}

const deleteDeps = "DELETE FROM deps WHERE left_id=? AND right_id=?"

// Delete breaks the association
func (s *SqliteStore) Delete(left, right string) error {
	leftID, err := s.nameToID(left)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}
	rightID, err := s.nameToID(right)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}
	_, err = s.db.Exec(deleteDeps, leftID, rightID)
	return err
}

// Close deinitializes the connection to the backend store
func (s *SqliteStore) Close() error {
	if !s.open {
		return fmt.Errorf("DB is alread closed")
	}
	return s.db.Close()
}
