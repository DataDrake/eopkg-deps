package main

import (
	"encoding/xml"
	"os"
    "text/template"
)

const digraph string = `digraph {
ranksep=2;
rankdir=LR;
{{range $package := .Packages}}{{$name := .Name}}{{range $dep := .RuntimeDependencies}}    "{{$name}}" -> "{{$dep}}";
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

type Package struct {
    Name string `xml:"Name"`
    RuntimeDependencies []string `xml:"RuntimeDependencies>Dependency"`
}

type Index struct {
    Packages []Package `xml:"Package"`
}

func main() {
	f, err := os.Open("/var/lib/eopkg/index/unstable/eopkg-index.xml")
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()

	d := xml.NewDecoder(f)
	var rawIndex Index
	err = d.Decode(&rawIndex)
	if err != nil {
		panic(err.Error())
	}
	err = outputTemplate.ExecuteTemplate(os.Stdout, "digraph", rawIndex)
    if err != nil {
        print(err.Error())
    }
}
