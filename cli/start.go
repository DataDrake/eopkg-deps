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

// Start marks a package for rebuilds
var Start = cmd.CMD{
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
func StartRun(r *cmd.RootCMD, c *cmd.CMD) {
	//flags := r.Flags.(*GlobalFlags)
	args := c.Args.(*StartArgs)
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
	err = s.StartToDo(args.Name)
	if err != nil {
		s.Close()
		fmt.Printf("Failed to mark for rebuilds , reason: '%s'\n", err.Error())
		os.Exit(1)
	}
    fmt.Printf("Successfully marked '%s' for rebuilds\n", args.Name)
	s.Close()
	os.Exit(0)
}
