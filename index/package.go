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

package index

// Package represents a single package and its immediate dependencies
type Package struct {
	Name     string `xml:"Name"`
	Releases []struct {
		Number int `xml:"release,attr"`
	} `xml:"History>Update"`
	RuntimeDependencies []struct {
		Name    string `xml:",chardata"`
		Release int    `xml:"releaseFrom,attr"`
	} `xml:"RuntimeDependencies>Dependency"`
}
