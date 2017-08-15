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
	"strings"

	"github.com/betterdoctor/duncan/chronos"
	"github.com/spf13/cobra"
)

var (
	schedule string
	follow   bool
	mem      int
	cpu      float64
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run a one-off process inside a remote container",
	Long: `The run command spins up a one-off remote container to execute
the supplied COMMAND. The COMMAND must exist inside the application's
Docker image or will result in failure. Logs of the running process are
printed to STDOUT.

duncan run -a APP -e ENV COMMAND

Example:

$ duncan run -a foo -e production rake stuff:junk

# to override default 1GB memory
$ duncan run -a foo -e production --mem 2 rake stuff:junk

If the command contains flags you will need to escape them
$ duncan run -a foo -e stage -- ls -lh

To run a scheduled job (PST)
$ duncan run -a foo -e stage --schedule R1/2017-08-15T00:00:00Z/PT30M some_task
`,
	Run: func(cmd *cobra.Command, args []string) {
		if app == "" {
			fmt.Println("must supply --app")
			os.Exit(1)
		}
		if env == "" {
			fmt.Println("must supply --env")
			os.Exit(1)
		}
		if len(args) == 0 {
			fmt.Println("must supply COMMAND to run")
			os.Exit(1)
		}
		command := strings.Join(args, " ")
		if err := chronos.RunCommand(app, env, command, schedule, cpu, mem, follow); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&app, "app", "a", "", "app to deploy")
	runCmd.Flags().StringVarP(&env, "env", "e", "", "deployment environment (stage, production)")
	runCmd.Flags().IntVarP(&mem, "mem", "", 1, "task memory in GB NOTE: must have available resources in cluster")
	runCmd.Flags().Float64VarP(&cpu, "cpu", "", 1.0, "task CPU NOTE: must have available resources in cluster")
	runCmd.Flags().StringVarP(&schedule, "schedule", "", "", "chronos job schedule (e.g., R1/2017-08-15T00:00:00Z/PT30M)")
	runCmd.Flags().BoolVarP(&follow, "follow", "f", false, "open link to task sandbox to follow streaming logs")
}
