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
)

func init() {
	cmd.Register(&Start)
}

// Start marks a package for rebuilds
var Start = cmd.Sub{
	Name:  "start",
	Alias: "to",
	Short: "Mark a package for rebuilds",
	Args:  &StartArgs{},
	Run:   StartRun,
}

// StartArgs contains the arguments for the "start" subcommand
type StartArgs struct {
	Name string `desc:"the name of the package to rebuild"`
}

// StartRun carries out the "start" subcommand
func StartRun(r *cmd.Root, c *cmd.Sub) {
	//flags := r.Flags.(*GlobalFlags)
	args := c.Args.(*StartArgs)
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
	if err = s.StartToDo(args.Name); err == sql.ErrNoRows {
		fmt.Printf("Package '%s' does not exist or you need to update\n", args.Name)
		os.Exit(1)
	}
	if err != nil {
		fmt.Printf("Failed to mark for rebuilds , reason: '%s'\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("Successfully marked '%s' for rebuilds\n", args.Name)
}
