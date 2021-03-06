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

package cli

import (
	"database/sql"
	"fmt"
	"github.com/DataDrake/cli-ng/v2/cmd"
	"github.com/DataDrake/eopkg-deps/storage"
	"os"
	"os/user"
	"sort"
	"text/tabwriter"
)

func init() {
	cmd.Register(&Forward)
}

// Forward gets a list of packages that this package depends on
var Forward = cmd.Sub{
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
	DependencyHeader = "Dependency\tSince Release\n"
	// DependencyHeaderColor is a table heading for forward dependencies, in color
	DependencyHeaderColor = "\033[1mDependency\tSince Release\n"
)

// ForwardRun carries out the "forward" subcommand
func ForwardRun(r *cmd.Root, c *cmd.Sub) {
	flags := r.Flags.(*GlobalFlags)
	args := c.Args.(*ForwardArgs)
	s := storage.NewStore()
	curr, err := user.Current()
	if err != nil {
		fmt.Printf(UserErrorFormat, err.Error())
		os.Exit(1)
	}
	if err = s.Open(curr.HomeDir + DefaultDBLocation); err != nil {
		fmt.Printf(DBOpenErrorFormat, err.Error())
		os.Exit(1)
	}
	defer s.Close()
	rights, err := s.GetForward(args.Package)
	if err == sql.ErrNoRows {
		fmt.Printf("Package '%s' does not exist or you need to update\n", args.Package)
		os.Exit(1)
	}
	if err != nil {
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
		fmt.Printf("No dependencies found.\n\n")
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	var rowFormat string
	if flags.NoColor {
		fmt.Fprintf(w, DependencyHeader)
		rowFormat = RowFormat
	} else {
		fmt.Fprintf(w, DependencyHeaderColor)
		rowFormat = RowFormatColor
	}
	for _, right := range rights {
		fmt.Fprintf(w, rowFormat, right.Name, right.Release)
	}
	w.Flush()
	fmt.Printf("\nTotal: %d\n", len(rights))
}
