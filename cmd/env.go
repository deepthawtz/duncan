// Copyright © 2020 Dylan Clendenin <dylan.clendenin@gmail.com>
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

	"github.com/deepthawtz/duncan/consul"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// envCmd represents the env command
var (
	envCmd = &cobra.Command{
		Use:   "env",
		Short: "Manage Consul key/values (ENV vars) for an app",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("must provide env subcommand, see: duncan env -h")
			os.Exit(1)
		},
	}

	envGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Display ENV vars for an app",
		Run: func(cmd *cobra.Command, args []string) {
			checkAppEnv(app, env)

			u := consul.EnvURL(app, env, true)
			env, err := consul.Read(u)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			printSorted(env)
		},
	}

	envSetCmd = &cobra.Command{
		Use:   "set KEY=VALUE [KEY2=VALUE2 ...]",
		Short: "Set one or more ENV var key/value pairs for an app",
		Run: func(cmd *cobra.Command, args []string) {
			checkAppEnv(app, env)
			validateKeyValues(args)

			u := consul.EnvURL(app, env, true)
			envVals, err := consul.Read(u)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			changes := make(map[string][2]string)
			for _, x := range args {
				parts := strings.Split(x, "=")
				k, v := parts[0], parts[1]
				prev, ok := envVals[k]
				if !ok {
					changes[k] = [2]string{"", v}
					continue
				}
				changes[k] = [2]string{prev, v}
			}

			if promptModifyEnvironment("set", "env", app, env, changes) {
				url := consul.TxnURL()
				env, err := consul.Write(app, env, url, args)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				printSorted(env)
			}
		},
	}

	envDelCmd = &cobra.Command{
		Use:   "del KEY [KEY ...]",
		Short: "Delete one or more ENV vars for an app",
		Run: func(cmd *cobra.Command, args []string) {
			checkAppEnv(app, env)
			validateKeys(args)

			u := consul.EnvURL(app, env, true)
			envVals, err := consul.Read(u)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			changes := make(map[string][2]string)
			for _, k := range args {
				prev, ok := envVals[k]
				if !ok {
					fmt.Printf("nothing to delete. no keys exists for %s\n", k)
					continue
				}
				changes[k] = [2]string{prev, ""}
			}

			if len(changes) == 0 {
				os.Exit(0)
			}

			if promptModifyEnvironment("delete", "env", app, env, changes) {
				url := consul.EnvURL(app, env, true)
				if err := consul.Delete(app, env, url, args); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
		},
	}
)

func init() {
	RootCmd.AddCommand(envCmd)
	envCmd.PersistentFlags().StringVarP(&app, "app", "a", "", "app to manage ENV vars for")
	envCmd.PersistentFlags().StringVarP(&env, "env", "e", "", "app environment (stage, production)")
	envSetCmd.Flags().BoolVarP(&force, "force", "f", false, "bypass prompt before setting env")
	envCmd.AddCommand(envSetCmd)
	envCmd.AddCommand(envGetCmd)
	envCmd.AddCommand(envDelCmd)
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

func promptModifyEnvironment(op, cmd, app, env string, changes map[string][2]string) bool {
	if force {
		return true
	}
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
	for k, transition := range changes {
		fmt.Printf("change %s from %s => %s\n", k, white(transition[0]), cyan(transition[1]))
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
