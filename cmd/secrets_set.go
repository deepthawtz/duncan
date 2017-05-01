// Copyright © 2016 Dylan Clendenin <dylan@betterdoctor.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
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

package cmd

import (
	"fmt"
	"os"

	"github.com/betterdoctor/duncan/vault"
	"github.com/spf13/cobra"
)

var secretsSetCmd = &cobra.Command{
	Use:   "set KEY=VALUE [KEY2=VALUE2 ...]",
	Short: "Set one or more secret key/value pairs for an app",
	Run: func(cmd *cobra.Command, args []string) {
		checkAppEnv(app, env)
		validateKeyValues(args)

		if promptModifyEnvironment("set", "secrets", app, env, args) {
			u := vault.SecretsURL(app, env)
			s, err := vault.Read(u)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			s, err = vault.Write(u, args, s)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}

			printSorted(s.KVPairs)
		}
	},
}

func init() {
	secretsCmd.AddCommand(secretsSetCmd)
}
