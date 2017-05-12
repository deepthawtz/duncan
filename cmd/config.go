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
	"regexp"
	"strings"
	"sync"

	"github.com/betterdoctor/duncan/consul"
	"github.com/betterdoctor/duncan/vault"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Search ENV/secrets across all applications",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("must call config subcommand")
			os.Exit(-1)
		},
	}

	configSearchCmd = &cobra.Command{
		Use:   "search PATTERN",
		Short: "Search ENV/secrets across applications by key",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("must provide PATTERN to search keys by")
				os.Exit(-1)
			}
			pattern := regexp.MustCompile(fmt.Sprintf(".*%s.*", args[0]))
			green := color.New(color.FgGreen, color.Bold).SprintFunc()
			matches := map[string]map[string]string{}
			apps := viper.GetStringSlice("apps")
			var (
				wg  sync.WaitGroup
				mux sync.Mutex
			)
			for _, app := range apps {
				if env == "" {
					env = "stage|production"
				}
				for _, e := range strings.Split(env, "|") {
					wg.Add(1)
					go func(app, env string) {
						defer wg.Done()
						ak := fmt.Sprintf("%s-%s", app, env)
						u := consul.EnvURL(app, env)
						// assume Consul read error means ACL restriction
						// or keyspace does not exist yet
						c, _ := consul.Read(u)
						if c != nil {
							for k, v := range c {
								if pattern.MatchString(k) {
									mux.Lock()
									if matches[ak] == nil {
										matches[ak] = make(map[string]string)
									}
									matches[ak][k] = v
									mux.Unlock()
								}
							}
						}
					}(app, e)

					wg.Add(1)
					go func(app, env string) {
						defer wg.Done()
						ak := fmt.Sprintf("%s-%s", app, env)
						u := vault.SecretsURL(app, env)
						s, _ := vault.Read(u)
						if s != nil {
							for k, v := range s.KVPairs {
								if pattern.MatchString(k) {
									mux.Lock()
									if matches[ak] == nil {
										matches[ak] = make(map[string]string)
									}
									matches[ak][k] = v
									mux.Unlock()
								}
							}
						}
					}(app, e)
				}
			}
			wg.Wait()
			for a, m := range matches {
				fmt.Println(green(a))
				printSorted(m)
			}
		},
	}
)

func init() {
	RootCmd.AddCommand(configCmd)
	configCmd.PersistentFlags().StringVarP(&env, "env", "e", "", "app environment (stage, production)")
	configCmd.AddCommand(configSearchCmd)
}
