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
	"github.com/DataDrake/eopkg-deps/index"
	"github.com/DataDrake/eopkg-deps/storage"
	"os"
	"os/user"
)

func init() {
	cmd.Register(&Update)
}

// Update creats a new datastore and populates it from the current eopkg index
var Update = cmd.Sub{
	Name:  "update",
	Alias: "up",
	Short: "Update rebuilds the datastore from the eopkg index",
	Run:   UpdateRun,
}

// UpdateRun carries out the "update" subcommand
func UpdateRun(r *cmd.Root, c *cmd.Sub) {
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
	if err = s.Open(curr.HomeDir + DefaultDBLocation); err != nil {
		fmt.Printf(DBOpenErrorFormat, err.Error())
		os.Exit(1)
	}
	defer s.Close()
	if err = s.Update(i); err != nil {
		fmt.Printf("Failed to update DB, reason: '%s'\n", err.Error())
		os.Exit(1)
	}
}
