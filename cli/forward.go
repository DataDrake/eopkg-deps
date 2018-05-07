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
	"os/user"
	"sort"
	"text/tabwriter"
)

// Forward gets a list of packages that this package depends on
var Forward = cmd.CMD{
	Name:  "forward",
	Alias: "fwd",
	Short: "Get this package's dependencies",
	Args:  &ForwardArgs{},
	Run:   ForwardRun,
}

// ForwardArgs contains the arguments for the "forward" subcommand
type ForwardArgs struct {
	Package string `desc:"the name of the package"`
}

const (
	// DependencyHeader is a table heading for forward dependencies
	DependencyHeader = "Dependency\tSince Release"
	// DependencyHeaderColor is a table heading for forward dependencies, in color
	DependencyHeaderColor = "\033[1mDependency\tSince Release"
)

// ForwardRun carries out the "forward" subcommand
func ForwardRun(r *cmd.RootCMD, c *cmd.CMD) {
	flags := r.Flags.(*GlobalFlags)
	args := c.Args.(*ForwardArgs)
	s := storage.NewStore()
	curr, err := user.Current()
	if err != nil {
		fmt.Printf(UserErrorFormat, err.Error())
		os.Exit(1)
	}
	err = s.Open(curr.HomeDir + DefaultDBLocation)
	if err != nil {
		fmt.Printf(DBOpenErrorFormat, err.Error())
		os.Exit(1)
	}
	rights, err := s.GetForward(args.Package)
	if err != nil {
		s.Close()
		fmt.Printf("Failed to get forward deps, reason: '%s'\n", err.Error())
		os.Exit(1)
	}
	sort.Sort(rights)
	if flags.NoColor {
		fmt.Printf(PackageFormat, args.Package)
	} else {
		fmt.Printf(PackageFormatColor, args.Package)
	}
	if len(rights) == 0 {
		s.Close()
		fmt.Println("No dependencies found.\n")
		os.Exit(0)
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	var rowFormat string
	if flags.NoColor {
		fmt.Println(DependencyHeader, args.Package)
		rowFormat = RowFormat
	} else {
		fmt.Println(DependencyHeaderColor, args.Package)
		rowFormat = RowFormatColor
	}
	for _, right := range rights {
		fmt.Fprintf(w, rowFormat, right.Name, right.Release)
	}
	w.Flush()
	fmt.Println()
	s.Close()
	os.Exit(0)
}
