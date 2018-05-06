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

import (
	"encoding/xml"
	"os"
	"text/template"
)

const digraph string = `digraph {
ranksep=2;
rankdir=LR;
{{range $package := .Packages}}{{$name := .Name}}{{range $dep := .RuntimeDependencies}}    "{{$name}}" -> "{{$dep.Name}}";
{{end}}{{end}}
}
`

var outputTemplate *template.Template

func init() {
	var err error
	outputTemplate = template.New("digraph")
	outputTemplate, err = outputTemplate.Parse(digraph)
	if err != nil {
		panic(err.Error())
	}
}

// Index represents all of the packages in the eopkg index
type Index struct {
	Packages []Package `xml:"Package"`
}

// NewIndex returns an uninitialized Index
func NewIndex() *Index {
	return &Index{}
}

// Load populates theis Index from an actual eopkg index
func (i *Index) Load(filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	d := xml.NewDecoder(f)
	return d.Decode(&i)
}

// Graph prints out a graph representation of an index
func (i *Index) Graph() error {
	return outputTemplate.ExecuteTemplate(os.Stdout, "digraph", i)
}
