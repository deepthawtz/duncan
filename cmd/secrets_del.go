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

// delCmd represents the del command
var delCmd = &cobra.Command{
	Use:   "del KEY [KEY ...]",
	Short: "Delete a secret key/value pair for an app",
	Run: func(cmd *cobra.Command, args []string) {
		checkAppEnv(app, env)
		validateKeys(args)

		for _, arg := range args {
			if err := vault.Delete(app, env, arg); err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
		}
	},
}

func init() {
	secretsCmd.AddCommand(delCmd)
}
