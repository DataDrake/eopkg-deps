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
)

// Root is the main command for this application
var Root *cmd.RootCMD

// GlobalFlags contains flags applicable to all sub-commands
type GlobalFlags struct {
	NoColor bool `short:"N" long:"no-color" desc:"Disable coloring of output text"`
}

func init() {
	// Build Application
	Root = &cmd.RootCMD{
		Name:  "eopkg-deps",
		Short: "Manage and work with eopkg dependencies",
		Flags: &GlobalFlags{
			NoColor: false,
		},
	}
	// Setup the Sub-Commands
	Root.RegisterCMD(&cmd.Help)
	Root.RegisterCMD(&Forward)
	Root.RegisterCMD(&Reverse)
	Root.RegisterCMD(&Rebuild)
	//Root.RegisterCMD(&Update)
}
