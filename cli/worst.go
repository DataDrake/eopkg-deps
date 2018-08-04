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
	"database/sql"
	"fmt"
	"github.com/DataDrake/cli-ng/cmd"
	"github.com/DataDrake/eopkg-deps/storage"
	"os"
	"os/user"
	"sort"
)

// Worst estimates what a full receursive rebuild would look like
var Worst = cmd.CMD{
	Name:  "worst",
	Alias: "ow",
	Short: "Calculate the worst-case rebuild list",
	Args:  &WorstArgs{},
	Run:   WorstRun,
}

// WorstArgs contains the arguments for the "worst" subcommand
type WorstArgs struct {
	Name string `desc:"the name of the package to rebuild"`
}

const (
	// WorstHeader is a table heading for required rebuilds
	WorstHeader = "Required Rebuilds"
	// WorstHeaderColor is a table heading for required rebuilds, in color
	WorstHeaderColor = "\033[1mRequired Rebuilds"
)

// WorstRun carries out the "worst" subcommand
func WorstRun(r *cmd.RootCMD, c *cmd.CMD) {
	flags := r.Flags.(*GlobalFlags)
	args := c.Args.(*WorstArgs)
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
	list, err := s.WorstToDo(args.Name)
	if err == sql.ErrNoRows {
		fmt.Printf("Package '%s' does not exist or you need to update\n", args.Name)
		s.Close()
		os.Exit(1)
	}
	if err != nil {
		s.Close()
		fmt.Printf("Failed to get todo list, reason: '%s'\n", err.Error())
		os.Exit(1)
	}
	if len(list) == 0 {
		s.Close()
		fmt.Println("No todo items found.\n")
		os.Exit(0)
	}
	sort.Sort(list)
	var rowFormat string
	if flags.NoColor {
		fmt.Println(WorstHeader)
		rowFormat = "%s\n"
	} else {
		fmt.Println(WorstHeaderColor)
		rowFormat = "\033[0m%s\n"
	}
	for _, item := range list {
		fmt.Printf(rowFormat, item.Name)
	}
	fmt.Println()
	if flags.NoColor {
		fmt.Printf("%s: %d\n", "Total", len(list))
	} else {
		fmt.Printf("\033[0m%s: %d\n", "Total", len(list))
	}
	s.Close()
	os.Exit(0)
}
