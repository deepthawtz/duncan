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

	"github.com/betterdoctor/duncan/consul"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var envSetCmd = &cobra.Command{
	Use:   "set KEY=VALUE [KEY2=VALUE2 ...]",
	Short: "Set one or more ENV var key/value pairs for an app",
	Run: func(cmd *cobra.Command, args []string) {
		checkAppEnv(app, env)
		validateKeyValues(args)

		if promptModifyEnvironment("set", "env", app, env, args) {
			host := viper.GetString("consul_host")
			token := viper.GetString("consul_token")
			url := fmt.Sprintf("https://%s/v1/txn?token=%s", host, token)
			env, err := consul.Write(app, env, url, args)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			printSorted(env)
		}
	},
}

func init() {
	envCmd.AddCommand(envSetCmd)
}
