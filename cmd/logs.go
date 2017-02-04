// Copyright © 2017 Dylan Clendenin <dylan@betterdoctor.com>
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

	"github.com/betterdoctor/duncan/logs"
	"github.com/spf13/cobra"
)

var utc bool

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Streams logs of your service",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if app == "" || env == "" {
			fmt.Println("must supply all flags for logs command")
			os.Exit(-1)
		}
		logs.Stream(app, env, utc)
	},
}

func init() {
	RootCmd.AddCommand(logsCmd)

	logsCmd.Flags().StringVarP(&app, "app", "a", "", "app to deploy")
	logsCmd.Flags().StringVarP(&env, "env", "e", "", "deployment environment (stage, production)")
	logsCmd.Flags().BoolVarP(&utc, "utc", "", false, "display log timestamps UTC (default: local time)")
}
