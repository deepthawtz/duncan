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

	"github.com/spf13/cobra"
)

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage Consul key/values (ENV vars) for an app",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("must call env subcommand")
		os.Exit(-1)
	},
}

func init() {
	RootCmd.AddCommand(envCmd)
	envCmd.PersistentFlags().StringVarP(&app, "app", "a", "", "app to manage ENV vars for")
	envCmd.PersistentFlags().StringVarP(&env, "env", "e", "", "app environment (stage, production)")
}
