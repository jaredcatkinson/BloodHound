/// 2>/dev/null ; gorun "$0" "$@" ; exit $?

// Copyright 2023 Specter Ops, Inc.
//
// Licensed under the Apache License, Version 2.0
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"

	"github.com/specterops/bloodhound/log"
	"github.com/specterops/bloodhound/packages/go/stbernard/command"
	"github.com/specterops/bloodhound/packages/go/stbernard/environment"
)

func main() {
	env := environment.NewEnvironment()
	var rawLvl = env[environment.LogLevelVarName]

	log.ConfigureDefaults()

	if rawLvl == "" {
		rawLvl = "warn"
	}

	if lvl, err := log.ParseLevel(rawLvl); err != nil {
		log.Errorf("Could not parse log level from %s: %v", environment.LogLevelVarName, err)
	} else {
		log.SetGlobalLevel(lvl)
	}

	if cmd, err := command.ParseCLI(env); errors.Is(err, command.ErrNoCmd) {
		log.Fatalf("No valid command specified")
	} else if errors.Is(err, command.ErrHelpRequested) {
		// No need to exit 1 if help was requested
		return
	} else if err != nil {
		log.Fatalf("Error while parsing command: %v", err)
	} else if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to run command `%s`: %v", cmd.Name(), err)
	} else {
		log.Infof("Command `%s` completed successfully", cmd.Name())
	}
}
