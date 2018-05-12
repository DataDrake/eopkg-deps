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
	"github.com/DataDrake/eopkg-deps/index"
)

// Store is a common interface for all kinds of backing store
type Store interface {
	// Open initializes a connection to the backend store
	Open(location string) error
	// GetForward returns: (left) -> *
	GetForward(lhs string) (Packages, error)
	// GetReverse returns: * -> (right)
	GetReverse(rhs string) (Packages, error)
	// GetToDo returns a list of unblocked packages that need to be rebuilt,
    // the count of remaining rebuilds, and the count of completed rebuilds
	GetToDo() (Packages, int, int, error)
    // StartToDo adds a new package to the todo list
    StartToDo(name string) error
	// Update clears the current store and rebuilds the contents from the provided index
	Update(i *index.Index) error
	// Close deinitializes the connection to the backend store
	Close() error
}

// NewStore gets a new version of the current preferred backing store
func NewStore() Store {
	return NewSqliteStore()
}
