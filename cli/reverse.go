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

package cli

import (
	"fmt"
	"github.com/DataDrake/cli-ng/cmd"
	"github.com/DataDrake/eopkg-deps/storage"
	"os"
    "sort"
)

// Reverse gets a list of packages that depend on this package
var Reverse = cmd.CMD{
	Name:  "reverse",
	Alias: "rev",
	Short: "Get this package's reverse dependencies",
	Args:  &ReverseArgs{},
	Run:   ReverseRun,
}

// ReverseArgs contains the arguments for the "reverse" subcommand
type ReverseArgs struct{
    Package string `desc:"the name of the package"`
}

// ReverseRun carries out the "Reverse" subcommand
func ReverseRun(r *cmd.RootCMD, c *cmd.CMD) {
	args := c.Args.(*ReverseArgs)
	s := storage.NewStore()
    err  := s.Open("/tmp/eopkg-deps.db")
	if err != nil {
		fmt.Printf("Failed to open DB, reason: '%s'\n", err.Error())
		os.Exit(1)
	}
	lefts, err := s.GetRight(args.Package)
	if err != nil {
		fmt.Printf("Failed to resolve reverse deps, reason: '%s'\n", err.Error())
		os.Exit(1)
	}
    sort.Strings(lefts)
    fmt.Printf("Name:\n\t%s\n\nReverse Dependencies:\n", args.Package)
    for _, left := range lefts {
        fmt.Printf("\t%s\n", left)
    }
    fmt.Println()
	os.Exit(0)
}
