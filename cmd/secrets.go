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
	"strings"

	"github.com/betterdoctor/duncan/vault"
	"github.com/spf13/cobra"
)

// secretsCmd represents the secrets command
var (
	secretsCmd = &cobra.Command{
		Use:   "secrets",
		Short: "Manage Vault secrets (ENV vars) for an app",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("must call secrets subcommand")
			os.Exit(-1)
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
				os.Exit(-1)
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

	secretsDelCmd = &cobra.Command{
		Use:   "del KEY [KEY ...]",
		Short: "Delete one or more secrets for an app",
		Run: func(cmd *cobra.Command, args []string) {
			checkAppEnv(app, env)
			validateKeys(args)

			if promptModifyEnvironment("delete", "secrets", app, env, args) {
				u := vault.SecretsURL(app, env)
				s, err := vault.Read(u)
				if err != nil {
					fmt.Println(err)
					os.Exit(-1)
				}
				if _, err := vault.Delete(u, args, s); err != nil {
					fmt.Println(err)
					os.Exit(-1)
				}
			}
		},
	}
)

func init() {
	RootCmd.AddCommand(secretsCmd)

	secretsCmd.PersistentFlags().StringVarP(&app, "app", "a", "", "app to manage secrets for")
	secretsCmd.PersistentFlags().StringVarP(&env, "env", "e", "", "app environment (stage, production)")
	secretsCmd.AddCommand(secretsGetCmd)
	secretsCmd.AddCommand(secretsSetCmd)
	secretsCmd.AddCommand(secretsDelCmd)
}

func checkAppEnv(app, env string) {
	if app == "" || env == "" {
		fmt.Println("must provide --app and --env flags")
		os.Exit(-1)
	}

	if env != "stage" && env != "production" {
		fmt.Printf("env %s is not a valid deployment environment\n", env)
		os.Exit(-1)
	}
}

func validateKeyValues(kvs []string) {
	if len(kvs) == 0 {
		fmt.Println("must provide key/value pairs in KEY=VALUE format")
		os.Exit(-1)
	}

	for _, k := range kvs {
		p := strings.Split(k, "=")
		// len should be at least 2 (edgecase w/ values that contain '=' character)
		if len(p) < 2 {
			fmt.Println("must provide key/value pairs in KEY=VALUE format")
			os.Exit(-1)
		}
	}
}

func validateKeys(keys []string) {
	if len(keys) == 0 {
		fmt.Println("must provide one or more keys")
		os.Exit(-1)
	}
	for _, k := range keys {
		p := strings.Split(k, "=")
		if len(p) > 1 {
			fmt.Println("KEY only must be provided, not KEY=VALUE")
			os.Exit(-1)
		}
	}
}
