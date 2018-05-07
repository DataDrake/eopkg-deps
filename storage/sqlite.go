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
	"fmt"
	"github.com/DataDrake/eopkg-deps/index"
	"github.com/jmoiron/sqlx"
	// Since this is the only place we will use sqlite directly
	_ "github.com/mattn/go-sqlite3"
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
CREATE TABLE IF NOT EXISTS packages (
    id   INTEGER PRIMARY KEY,
    name TEXT,
    rel  INTEGER
);

CREATE TABLE IF NOT EXISTS deps (
    left_id   INTEGER,
    right_id  INTEGER,
    rel INTEGER
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

const getPackage = "SELECT id FROM packages WHERE name=?"

func (s *SqliteStore) nameToID(name string) (int, error) {
	var id int
	row := s.db.QueryRowx(getPackage, name)
	err := row.Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}

const getRHS = `
SELECT name, rel2 AS rel FROM packages INNER JOIN (
    SELECT right_id, rel AS rel2 FROM deps WHERE left_id=?
) ON packages.id=right_id
`

// GetForward returns: (lhs) -> *
func (s *SqliteStore) GetForward(lhs string) (Packages, error) {
	lhsID, err := s.nameToID(lhs)
	if err != nil {
		return nil, err
	}
	rhs := make(Packages, 0)
	rows, err := s.db.Queryx(getRHS, lhsID)
	if err != nil {
		return rhs, err
	}
	for rows.Next() {
		var p Package
		err := rows.StructScan(&p)
		if err != nil {
			return rhs, err
		}
		rhs = append(rhs, p)
	}
	return rhs, err
}

const getLHS = `
SELECT name, rel2 AS rel FROM packages INNER JOIN (
    SELECT left_id, rel AS rel2 FROM deps WHERE right_id=?
) ON packages.id=left_id
`

// GetReverse returns: * -> (rhs)
func (s *SqliteStore) GetReverse(rhs string) (Packages, error) {
	rhsID, err := s.nameToID(rhs)
	if err != nil {
		return nil, err
	}
	lhs := make(Packages, 0)
	rows, err := s.db.Queryx(getLHS, rhsID)
	if err != nil {
		return lhs, err
	}
	for rows.Next() {
		var p Package
		err := rows.StructScan(&p)
		if err != nil {
			return lhs, err
		}
		lhs = append(lhs, p)
	}
	return lhs, err
}

const dropTables = `
    DROP TABLE IF EXISTS packages;
    DROP TABLE IF EXISTS deps
`

const insertPackage = "INSERT INTO packages VALUES (?,?,?)"
const insertDep = "INSERT INTO deps VALUES (?,?,?)"

// Update rebuilds the store from an Index
func (s *SqliteStore) Update(i *index.Index) error {
	s.db.MustExec(dropTables)
	s.createTables()
	tx := s.db.MustBegin()
	pkgStmt, err := tx.Preparex(insertPackage)
	if err != nil {
		return err
	}
	// Get ID mappings
	idMap := make(map[string]int)
	for id, pkg := range i.Packages {
		idMap[pkg.Name] = id
		_, err = pkgStmt.Exec(id, pkg.Name, pkg.Releases[0].Number)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	depStmt, err := tx.Preparex(insertDep)
	if err != nil {
		tx.Rollback()
		return err
	}
	//Get left-right mappings
	for leftID, lpkg := range i.Packages {
		for _, rpkg := range lpkg.RuntimeDependencies {
			rightID := idMap[rpkg.Name]
			_, err = depStmt.Exec(leftID, rightID, rpkg.Release)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}

// Close deinitializes the connection to the backend store
func (s *SqliteStore) Close() error {
	if !s.open {
		return fmt.Errorf("DB is alread closed")
	}
	return s.db.Close()
}
