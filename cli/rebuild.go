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
	"github.com/DataDrake/eopkg-deps/index"
	"github.com/DataDrake/eopkg-deps/storage"
	"os"
	"os/user"
)

// Rebuild creats a new datastore and populates it from the currrent eopkg index
var Rebuild = cmd.CMD{
	Name:  "rebuild",
	Alias: "rbd",
	Short: "Rebuilds the datastore from the eopkg index",
	Args:  &RebuildArgs{},
	Run:   RebuildRun,
}

// RebuildArgs contains the arguments for the "rebuild" subcommand
type RebuildArgs struct{}

// RebuildRun carries out the "rebuild" subcommand
func RebuildRun(r *cmd.RootCMD, c *cmd.CMD) {
	//args := c.Args.(*RebuildArgs)
	i := index.NewIndex()
	err := i.Load("/var/lib/eopkg/index/Unstable/eopkg-index.xml")
	if err != nil {
		fmt.Printf("Failed to load index, reason: '%s'\n", err.Error())
		os.Exit(1)
	}
	//i.Graph()
	s := storage.NewStore()
	curr, err := user.Current()
	if err != nil {
		fmt.Printf("Failed to get user, reason: '%s'\n", err.Error())
		os.Exit(1)
	}
	err = s.Open(curr.HomeDir + "/.cache/eopkg-deps.db")
	if err != nil {
		fmt.Printf("Failed to open DB, reason: '%s'\n", err.Error())
		os.Exit(1)
	}
	err = s.Rebuild(i)
	if err != nil {
		fmt.Printf("Failed to rebuild DB, reason: '%s'\n", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
