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

// Update creats a new datastore and populates it from the current eopkg index
var Update = cmd.CMD{
	Name:  "update",
	Alias: "up",
	Short: "Update rebuilds the datastore from the eopkg index",
	Args:  &UpdateArgs{},
	Run:   UpdateRun,
}

// UpdateArgs contains the arguments for the "update" subcommand
type UpdateArgs struct{}

// UpdateRun carries out the "update" subcommand
func UpdateRun(r *cmd.RootCMD, c *cmd.CMD) {
	//args := c.Args.(*RebuildArgs)
	i := index.NewIndex()
	err := i.Load(DefaultIndexLocation)
	if err != nil {
		fmt.Printf("Failed to load index, reason: '%s'\n", err.Error())
		os.Exit(1)
	}
	//i.Graph()
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
	err = s.Update(i)
	if err != nil {
		fmt.Printf("Failed to update DB, reason: '%s'\n", err.Error())
		s.Close()
		os.Exit(1)
	}
	s.Close()
	os.Exit(0)
}
