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
	cmd.Register(&Reverse)
}

// Reverse gets a list of packages that depend on this package
var Reverse = cmd.Sub{
	Name:  "reverse",
	Alias: "rev",
	Short: "Get this package's reverse dependencies",
	Args:  &ReverseArgs{},
	Run:   ReverseRun,
}

// ReverseArgs contains the arguments for the "reverse" subcommand
type ReverseArgs struct {
	Package string `desc:"the name of the package"`
}

const (
	// ReverseDependencyHeader is a table heading for reverse dependencies
	ReverseDependencyHeader = "Reverse Dependency\tRelease\n"
	// ReverseDependencyHeaderColor is a table heading for reverse dependencies, in color
	ReverseDependencyHeaderColor = "\033[1mReverse Dependency\tRelease\n"
)

// ReverseRun carries out the "Reverse" subcommand
func ReverseRun(r *cmd.Root, c *cmd.Sub) {
	flags := r.Flags.(*GlobalFlags)
	args := c.Args.(*ReverseArgs)
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
	lefts, err := s.GetReverse(args.Package)
	if err == sql.ErrNoRows {
		fmt.Printf("Package '%s' does not exist or you need to update\n", args.Package)
		os.Exit(1)
	}
	if err != nil {
		fmt.Printf("Failed to resolve reverse deps, reason: '%s'\n", err.Error())
		os.Exit(1)
	}
	sort.Sort(lefts)
	if flags.NoColor {
		fmt.Printf(PackageFormat, args.Package)
	} else {
		fmt.Printf(PackageFormatColor, args.Package)
	}

	if len(lefts) == 0 {
		fmt.Printf("No reverse dependencies found.\n\n")
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	var rowFormat string
	if flags.NoColor {
		fmt.Fprintf(w, ReverseDependencyHeader)
		rowFormat = RowFormat
	} else {
		fmt.Fprintf(w, ReverseDependencyHeaderColor)
		rowFormat = RowFormatColor
	}
	for _, left := range lefts {
		fmt.Fprintf(w, rowFormat, left.Name, left.Release)
	}
	w.Flush()
	fmt.Printf("\nTotal: %d\n", len(lefts))
}
