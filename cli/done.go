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
	"strings"
)

// Done marks a package as rebuilt and optionally marks its reverse dependencies for rebuilds
var Done = cmd.CMD{
	Name:  "done",
	Alias: "do",
	Short: "Mark a package as rebuilt, marking reverse deps for rebuilds",
	Args:  &DoneArgs{},
	Run:   DoneRun,
}

// DoneArgs contains the arguments for the "done" subcommand
type DoneArgs struct {
	Name     string `desc:"the name of the package that was rebuilt"`
	Continue string `desc:"queue the reverse deps? (yes/no)"`
}

// DoneRun carries out the "done" subcommand
func DoneRun(r *cmd.RootCMD, c *cmd.CMD) {
	//flags := r.Flags.(*GlobalFlags)
	args := c.Args.(*DoneArgs)
	var Continue bool
	args.Continue = strings.ToLower(args.Continue)
	switch args.Continue {
	case "yes", "y", "true", "t":
		Continue = true
	case "no", "n", "false", "f":
		Continue = false
	default:
		fmt.Println("Coninue must be some flavor of (Y)es or (N)o")
		os.Exit(1)
	}
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
	err = s.DoneToDo(args.Name, Continue)
	if err != nil {
		s.Close()
		fmt.Printf("Failed to mark as rebuilt , reason: '%s'\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("Successfully marked '%s' as rebuilt\n", args.Name)
	s.Close()
	os.Exit(0)
}
