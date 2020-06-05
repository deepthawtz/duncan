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
	"fmt"
	"os"
	"strings"

	"github.com/deepthawtz/duncan/vault"
	"github.com/spf13/cobra"
)

// secretsCmd represents the secrets command
var (
	secretsCmd = &cobra.Command{
		Use:   "secrets",
		Short: "Manage Vault secrets (ENV vars) for an app",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("must call secrets subcommand")
			os.Exit(1)
		},
	}

	secretsGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Display secrets for an app",
		Run: func(cmd *cobra.Command, args []string) {
			checkAppEnv(app, env)

			u := vault.SecretsURL(app, env)
			s, err := vault.Read(u)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			printSorted(s.KVPairs)
		},
	}

	secretsSetCmd = &cobra.Command{
		Use:   "set KEY=VALUE [KEY2=VALUE2 ...]",
		Short: "Set one or more secret key/value pairs for an app",
		Run: func(cmd *cobra.Command, args []string) {
			checkAppEnv(app, env)
			validateKeyValues(args)

			u := vault.SecretsURL(app, env)
			secrets, err := vault.Read(u)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			changes := make(map[string][2]string)
			for _, x := range args {
				parts := strings.Split(x, "=")
				k, v := parts[0], parts[1]
				prev, ok := secrets.KVPairs[k]
				if !ok {
					changes[k] = [2]string{"", v}
					continue
				}
				changes[k] = [2]string{prev, v}
			}

			if promptModifyEnvironment("set", "secrets", app, env, changes) {
				s, err := vault.Write(u, args, secrets)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				printSorted(s.KVPairs)
			}
		},
	}

	secretsDelCmd = &cobra.Command{
		Use:   "del KEY [KEY ...]",
		Short: "Delete one or more secrets for an app",
		Run: func(cmd *cobra.Command, args []string) {
			checkAppEnv(app, env)
			validateKeys(args)

			u := vault.SecretsURL(app, env)
			secrets, err := vault.Read(u)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			changes := make(map[string][2]string)
			for _, k := range args {
				prev, ok := secrets.KVPairs[k]
				if !ok {
					fmt.Printf("nothing to delete. no keys exists for %s\n", k)
					continue
				}
				changes[k] = [2]string{prev, ""}
			}

			if len(changes) == 0 {
				os.Exit(0)
			}

			if promptModifyEnvironment("delete", "secrets", app, env, changes) {
				if _, err := vault.Delete(u, args, secrets); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
		},
	}
)

func init() {
	RootCmd.AddCommand(secretsCmd)

	secretsCmd.PersistentFlags().StringVarP(&app, "app", "a", "", "app to manage secrets for")
	secretsCmd.PersistentFlags().StringVarP(&env, "env", "e", "", "app environment (stage, production)")
	secretsSetCmd.Flags().BoolVarP(&force, "force", "f", false, "bypass prompt before setting env")
	secretsCmd.AddCommand(secretsGetCmd)
	secretsCmd.AddCommand(secretsSetCmd)
	secretsCmd.AddCommand(secretsDelCmd)
}

func checkAppEnv(app, env string) {
	if app == "" || env == "" {
		fmt.Println("must provide --app and --env flags")
		os.Exit(1)
	}
}

func validateKeyValues(kvs []string) {
	if len(kvs) == 0 {
		fmt.Println("must provide key/value pairs in KEY=VALUE format")
		os.Exit(1)
	}

	for _, k := range kvs {
		p := strings.Split(k, "=")
		// len should be at least 2 (edgecase w/ values that contain '=' character)
		if len(p) < 2 {
			fmt.Println("must provide key/value pairs in KEY=VALUE format")
			os.Exit(1)
		}
	}
}

func validateKeys(keys []string) {
	if len(keys) == 0 {
		fmt.Println("must provide one or more keys")
		os.Exit(1)
	}
	for _, k := range keys {
		p := strings.Split(k, "=")
		if len(p) > 1 {
			fmt.Println("KEY only must be provided, not KEY=VALUE")
			os.Exit(1)
		}
	}
}
