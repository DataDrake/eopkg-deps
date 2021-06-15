//
// Copyright 2018-2021 Bryan T. Meyers <root@datadrake.com>
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

// Package is Storage's Representation of a Package
type Package struct {
	Name    string `db:"name"`
	Release int    `db:"rel"`
}

// Packages is a sortable type for a list of Package struct
type Packages []Package

// Len returns the length of the list
func (pkgs Packages) Len() int {
	return len(pkgs)
}

// Less is based on the names of the packages only
func (pkgs Packages) Less(i, j int) bool {
	return pkgs[i].Name < pkgs[j].Name
}

// Swap switches the packages when sortint
func (pkgs Packages) Swap(i, j int) {
	pkgs[i], pkgs[j] = pkgs[j], pkgs[i]
}
