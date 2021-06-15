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
	"fmt"
	"github.com/DataDrake/cli-ng/v2/cmd"
	"github.com/DataDrake/eopkg-deps/storage"
	"os"
	"os/user"
)

func init() {
	cmd.Register(&Reset)
}

// Reset clears the current todo
var Reset = cmd.Sub{
	Name:  "reset",
	Alias: "clr",
	Short: "Clear the entire todo list",
	Run:   ResetRun,
}

// ResetRun carries out the "reset" subcommand
func ResetRun(r *cmd.Root, c *cmd.Sub) {
	//flags := r.Flags.(*GlobalFlags)
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
	if err = s.ResetToDo(); err != nil {
		fmt.Printf("Failed to reset ToDo list , reason: '%s'\n", err.Error())
		os.Exit(1)
	}
	fmt.Println("Successfully marked reset ToDo list")
}
