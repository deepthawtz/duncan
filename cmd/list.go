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

	"github.com/deepthawtz/duncan/k8s"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List applications",
	Run: func(cmd *cobra.Command, args []string) {
		k8sClient, err := k8s.NewClient()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err := k8sClient.List(app, env); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&app, "app", "a", "", "optionally filter by app")
	listCmd.Flags().StringVarP(&env, "env", "e", "", "optionally filter by environment (stage, production)")
}
