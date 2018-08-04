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
)

// Reset clears the current todo
var Reset = cmd.CMD{
	Name:  "reset",
	Alias: "clr",
	Short: "Clear the entire todo list",
	Args:  &ResetArgs{},
	Run:   ResetRun,
}

// ResetArgs contains the arguments for the "reset" subcommand
type ResetArgs struct{}

// ResetRun carries out the "reset" subcommand
func ResetRun(r *cmd.RootCMD, c *cmd.CMD) {
	//flags := r.Flags.(*GlobalFlags)
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
	err = s.ResetToDo()
	if err != nil {
		s.Close()
		fmt.Printf("Failed to reset ToDo list , reason: '%s'\n", err.Error())
		os.Exit(1)
	}
	fmt.Println("Successfully marked reset ToDo list")
	s.Close()
	os.Exit(0)
}
