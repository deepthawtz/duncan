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
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
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

func printSorted(m map[string]string) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s=%s\n", k, m[k])
	}
}

func promptModifyEnvironment(op, cmd, app, env string, args []string) bool {
	white := color.New(color.FgWhite, color.Bold).SprintFunc()
	red := color.New(color.FgRed, color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
	fmt.Printf("You are about to modify the the following envrionment:\n\n")
	fmt.Printf(white("  app: %s\n"), yellow(app))
	if env == "production" {
		env = red(env)
	} else {
		env = green(env)
	}
	fmt.Printf(white("  env: %s\n"), env)
	fmt.Printf(white("  command: %s\n"), cyan(cmd))
	fmt.Printf(white("  operation: %s\n"), cyan(op))
	if cmd == "secrets" && op == "set" {
		fmt.Printf(white("  encryption: %s\n"), green("enabled"))
	}
	fmt.Println()
	for _, el := range args {
		fmt.Println(el)
	}
	if op == "set" && cmd == "env" {
		fmt.Printf("\n%s ", red("WARNING:"))
		fmt.Printf(white("environment variables set w/ env command are NOT encrypted\n"))
		fmt.Println(white("         this command should not be used to store sensitive values like passwords or tokens"))
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(white("\nare you sure? (yes/no): "))
	resp, _ := reader.ReadString('\n')

	resp = strings.TrimSpace(resp)
	if resp != "yes" {
		return false
	}
	return true
}
