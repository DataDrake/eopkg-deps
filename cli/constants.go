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

// Paths
const (
	DefaultDBLocation    = "/.cache/eopkg-deps.db"
	DefaultIndexLocation = "/var/lib/eopkg/index/Unstable/eopkg-index.xml"
)

// Error Strings
const (
	DBOpenErrorFormat = "Failed to open DB, reason: '%s'\n"
	UserErrorFormat   = "Failed to get user, reason: '%s'\n"
)

// Format Strings
const (
	PackageFormat      = "Package: %s\n\n"
	PackageFormatColor = "\033[1mPackage:\033[0m %s\n\n"
	RowFormat          = "%s\t%d\n"
	RowFormatColor     = "\033[0m%s\t%d\n"
)
