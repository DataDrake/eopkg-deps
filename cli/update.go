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
	"github.com/DataDrake/cli-ng/cmd"
	//"github.com/DataDrake/eopkg-deps/storage"
	"os"
)

// Update opens the datastore and updates it based on the current eopkg index
var Update = cmd.CMD{
	Name:  "update",
	Alias: "up",
	Short: "Updates the datastore from the eopkg index",
	Args:  &UpdateArgs{},
	Run:   UpdateRun,
}

// UpdateArgs contains the arguments for the "update" subcommand
type UpdateArgs struct{}

// UpdateRun carries out the "update" subcommand
func UpdateRun(r *cmd.RootCMD, c *cmd.CMD) {
	//args := c.Args.(*UpdateArgs)
	os.Exit(0)
}
