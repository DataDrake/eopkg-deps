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
	"github.com/DataDrake/eopkg-deps/index"
	"github.com/jmoiron/sqlx"
	"strings"
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
    rel       INTEGER
);

CREATE TABLE IF NOT EXISTS todo (
    name       TEXT,
    package_id INTEGER,
    done       BOOLEAN
);
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

const getToDo = "SELECT count(*) FROM todo WHERE name=? AND done=FALSE"
const insertToDo = "INSERT OR REPLACE INTO todo VALUES (?, ?, FALSE)"

// StartToDo adds a new package to the todo list
func (s *SqliteStore) StartToDo(name string) error {
	var count int
	err := s.db.Get(&count, getToDo, name)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("Rebuild for package '%s' has already started", name)
	}
	id, err := s.nameToID(name)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(insertToDo, name, id)
	return err
}

const checkDone = "SELECT done FROM todo WHERE name=? AND done=false"
const markDone = "UPDATE todo SET done=TRUE WHERE name=?"

const insertReverse = `
INSERT INTO todo
    SELECT name, id, FALSE FROM packages INNER JOIN (
        SELECT left_id FROM deps WHERE right_id=?
    ) ON packages.id=left_id
    WHERE id NOT IN (SELECT package_id FROM todo)
`

// DoneToDo marks a package as complete and optionally queues its reverse deps
func (s *SqliteStore) DoneToDo(name string, Continue bool) error {
    done := false
	err := s.db.Get(&done, checkDone, name)
    if err == sql.ErrNoRows {
        return fmt.Errorf("Package '%s' is not in the todo list", name)
    }
	if err != nil {
		return err
	}
    if done {
        return fmt.Errorf("Package '%s' is already marked 'Done'", name)
    }
	_, err = s.db.Exec(markDone, name)
	if err != nil {
		return err
	}
	if Continue {
		id, err := s.nameToID(name)
		_, err = s.db.Exec(insertReverse, id)
		if err != nil {
			return err
		}
	}
	return err
}

const getUnblocked = `
SELECT name FROM todo WHERE package_id
NOT IN (
    SELECT left_id FROM deps INNER JOIN (
        SELECT package_id FROM todo WHERE done=FALSE
    ) ON right_id=package_id
) AND done=FALSE;
`
const getToDoCount = `SELECT count(*) FROM todo WHERE done=FALSE`
const getToDoDone = `SELECT count(*) FROM todo WHERE done=TRUE`

// GetToDo gets the currently unblocked packages to rebuild
func (s *SqliteStore) GetToDo() (Packages, int, int, error) {
	unblocked := make(Packages, 0)
	rows, err := s.db.Queryx(getUnblocked)
	if err != nil {
		return unblocked, 0, 0, err
	}
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return unblocked, 0, 0, err
		}
		unblocked = append(unblocked, Package{name, 0})
	}
	var count int
	err = s.db.Get(&count, getToDoCount)
	if err != nil {
		return unblocked, 0, 0, err
	}
	var done int
	err = s.db.Get(&done, getToDoDone)
	if err != nil {
		return unblocked, 0, 0, err
	}
	return unblocked, count, done, err
}

const getWorst = `
WITH RECURSIVE traverse AS (
    SELECT left_id FROM deps INNER JOIN packages
    ON right_id=id WHERE name=?
    UNION ALL
    SELECT deps.left_id FROM deps
        INNER JOIN traverse
        ON deps.right_id=traverse.left_id
)
SELECT name FROM traverse INNER JOIN packages
ON id=left_id GROUP BY name;
`

// WorstToDo gets a worst-case list of packages to rebuild
func (s *SqliteStore) WorstToDo(name string) (Packages, error) {
	list := make(Packages, 0)
	rows, err := s.db.Queryx(getWorst, name)
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var pName string
		err := rows.Scan(&pName)
		if err != nil {
			return list, err
		}
		list = append(list, Package{pName, 0})
	}
	return list, err
}

const resetToDo = "DELETE FROM todo"

// ResetToDo clears the todo list
func (s *SqliteStore) ResetToDo() error {
	_, err := s.db.Exec(resetToDo)
	return err
}

const dropTables = `
    DROP TABLE IF EXISTS packages;
    DROP TABLE IF EXISTS deps;
    DROP TABLE IF EXISTS todo;
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
		// skip -devel and -dbginfo packages
		if strings.HasSuffix(pkg.Name, "-dbginfo") || strings.HasSuffix(pkg.Name, "-devel") {
			continue
		}
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
