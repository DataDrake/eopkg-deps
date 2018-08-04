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
)

// ToDo gets a list of packages that still need to be rebuilt
var ToDo = cmd.CMD{
	Name:  "todo",
	Alias: "td",
	Short: "Get packages to rebuild",
	Args:  &ToDoArgs{},
	Run:   ToDoRun,
}

// ToDoArgs contains the arguments for the "todo" subcommand
type ToDoArgs struct{}

const (
	// ToDoHeader is a table heading for remaining packages
	ToDoHeader = "Unblocked Packages"
	// ToDoHeaderColor is a table heading for remaining packages, in color
	ToDoHeaderColor = "\033[1mUnblocked Packages"
)

// ToDoRun carries out the "todo" subcommand
func ToDoRun(r *cmd.RootCMD, c *cmd.CMD) {
	flags := r.Flags.(*GlobalFlags)
	//args := c.Args.(*ToDoArgs)
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
	var rowFormat string
	if flags.NoColor {
		rowFormat = "%s\n"
	} else {
		rowFormat = "\033[0m%s\n"
	}
	unblocked, count, done, err := s.GetToDo()
	if err != nil {
		s.Close()
		fmt.Printf("Failed to get todo list, reason: '%s'\n", err.Error())
		os.Exit(1)
	}
	sort.Sort(unblocked)
	if len(unblocked) == 0 {
		s.Close()
		fmt.Println("No todo items found.")
		goto DONE
	}
	if flags.NoColor {
		fmt.Println(ToDoHeader)
	} else {
		fmt.Println(ToDoHeaderColor)
	}
	for _, item := range unblocked {
		fmt.Printf(rowFormat, item.Name)
	}
DONE:
	fmt.Println()
	if flags.NoColor {
		fmt.Printf("%-10s: %d\n", "Unblocked", len(unblocked))
		fmt.Printf("%-10s: %d\n", "Queued", count)
		fmt.Printf("%-10s: %d\n", "Completed", done)
	} else {
		fmt.Printf("\033[0m%-10s: %d\n", "Unblocked", len(unblocked))
		fmt.Printf("\033[0m%-10s: %d\n", "Queued", count)
		fmt.Printf("\033[0m%-10s: %d\n", "Completed", done)
	}
	fmt.Println()
	s.Close()
	os.Exit(0)
}
