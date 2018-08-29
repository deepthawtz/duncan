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
	"os/user"
	"strings"

	"github.com/betterdoctor/duncan/deployment"
	"github.com/betterdoctor/duncan/docker"
	"github.com/betterdoctor/duncan/k8s"
	"github.com/betterdoctor/duncan/marathon"
	"github.com/betterdoctor/kit/notify"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	app, env, tag, repo, prev string
	force                     bool
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy an application",
	Long: `Deploy an application by specified tag.

Example:

$ duncan deploy --app APP --env ENV --tag TAG [--repo DOCKER_REPO]

NOTE: tag must exist in docker registry
`,

	Run: func(cmd *cobra.Command, args []string) {
		validateDeployFlags()

		var err error

		if viper.GetString("kubernetes_host") != "" {
			prev, err = k8s.CurrentTag(app, env, repo)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			prev, err = marathon.CurrentTag(app, env, repo)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		if promptDeploy() {
			if viper.GetString("kubernetes_host") != "" {
				if err := k8s.Deploy(app, env, tag, repo); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			} else {
				if err := marathon.Deploy(app, env, tag, repo); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}

			diff := "redeployed"
			if tag != prev {
				diff = deployment.GithubDiffLink(repo, prev, tag)
			}
			fmt.Println(diff)

			username := "bot"
			u, err := user.Current()
			if err == nil {
				username = u.Username
			}
			if err := notify.Slack(
				viper.GetString("slack_webhook_url"),
				fmt.Sprintf("%s %s (%s)", app, env, tag),
				fmt.Sprintf("%s :shipit: *%s %s* deployed by %s (diff: %s)", emoji(env), app, env, username, diff),
			); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)
	deployCmd.Flags().StringVarP(&app, "app", "a", "", "app to deploy")
	deployCmd.Flags().StringVarP(&env, "env", "e", "", "deployment environment (stage, production)")
	deployCmd.Flags().StringVarP(&tag, "tag", "t", "", "tag to deploy")
	deployCmd.Flags().StringVarP(&repo, "repo", "r", "", "(optional) if docker repo/image name differs from app name")
	deployCmd.Flags().BoolVarP(&force, "force", "f", false, "bypass prompt before deploying")
}

func validateDeployFlags() {
	if app == "" || env == "" || tag == "" {
		fmt.Println("must supply all flags for deploy command")
		os.Exit(1)
	}

	if env != "stage" && env != "production" {
		fmt.Printf("env %s is not a valid deployment environment\n", env)
		os.Exit(1)
	}

	// if no --repo flag use app name as repo name
	// this is important to use
	if repo == "" {
		repo = app
	}

	if err := docker.VerifyTagExists(repo, tag); err != nil {
		prefix := viper.GetString("docker_repo_prefix")
		fmt.Printf("could not verify %s/%s:%s exists: %s\n", prefix, repo, tag, err)
		os.Exit(1)
	}
}

func emoji(env string) string {
	if env == "production" {
		return ":balloon:"
	}

	return ""
}

func promptDeploy() bool {
	if force {
		return true
	}
	white := color.New(color.FgWhite, color.Bold).SprintFunc()
	red := color.New(color.FgRed, color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
	fmt.Printf("You are about to deploy:\n\n")
	fmt.Printf(white("  app: %s\n"), yellow(app))
	if env == "production" {
		fmt.Printf(white("  env: %s\n"), red(env))
	} else {
		fmt.Printf(white("  env: %s\n"), green(env))
	}
	fmt.Printf(white("  tag: %s => %s\n"), white(prev), cyan(tag))

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(white("\nare you sure? (yes/no): "))
	resp, _ := reader.ReadString('\n')

	resp = strings.TrimSpace(resp)
	if resp != "yes" {
		fmt.Println("phew... that was close")
		return false
	}
	return true
}
