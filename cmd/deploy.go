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

	"github.com/betterdoctor/duncan/chronos"
	"github.com/betterdoctor/duncan/docker"
	"github.com/betterdoctor/duncan/marathon"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	app, env, tag, marathonPath string
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy a docker app",
	Long: `For example:

duncan deploy --app dex --env stage --tag release-1.2.3`,

	Run: func(cmd *cobra.Command, args []string) {
		sanityCheck()
		if err := marathon.Deploy(app, env, tag); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		if err := chronos.Deploy(app, env, tag); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)
	deployCmd.Flags().StringVarP(&app, "app", "a", "", "app to deploy")
	deployCmd.Flags().StringVarP(&env, "env", "e", "", "deployment environment (stage, production)")
	deployCmd.Flags().StringVarP(&tag, "tag", "t", "", "git tag to deploy")
}

// sanityCheck checks that the most important pre-conditions are met before
// a deployment can proceed and dies if anything is not right
func sanityCheck() {
	if app == "" || env == "" || tag == "" {
		fmt.Println("must supply all flags for deploy command")
		os.Exit(-1)
	}

	if !appExists() {
		fmt.Printf("app %s does not exist yet\n", app)
		os.Exit(-1)
	}

	if env != "stage" && env != "production" {
		fmt.Printf("env %s is not a valid deployment environment\n", env)
		os.Exit(-1)
	}

	marathonPath = viper.GetString("marathon_json_path")
	if !marathonPathExists() {
		fmt.Printf("marathon path %s does not exist\n", marathonPath)
		os.Exit(-1)
	}

	if !docker.TagExists(app, tag) {
		fmt.Printf("docker tag %s does not exist for %s\n", tag, app)
		os.Exit(-1)
	}
}

func appExists() bool {
	apps := viper.GetStringMap("apps")
	for a := range apps {
		if a == app {
			return true
		}
	}
	return false
}

func marathonPathExists() bool {
	_, err := os.Stat(marathonPath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}
