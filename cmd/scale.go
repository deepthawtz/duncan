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
	"strconv"
	"strings"

	"github.com/betterdoctor/duncan/autoscaling"
	"github.com/betterdoctor/duncan/deployment"
	"github.com/betterdoctor/duncan/marathon"
	"github.com/betterdoctor/kit/notify"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type scaleEvent struct {
	App   string
	Env   string
	Procs []*proc
}

type proc struct {
	InstanceType string
	Previous     int
	Current      int
}

// scaleCmd represents the scale command
var scaleCmd = &cobra.Command{
	Use:   "scale",
	Short: "Scale an app process",
	Long: `Scale processes within an application group by name and count

Examples:

duncan scale web=2 --app foo --env production
duncan scale web=2 worker=5 --app foo --env production

If application cannot scale due to insufficient cluster resources an error will be returned
	`,
	Run: func(cmd *cobra.Command, args []string) {
		checkAppEnv(app, env)
		rules, err := parseScaleRules(args)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			fmt.Println("USAGE: duncan scale web=3 [worker=2, ...] --app foo --env production")
			os.Exit(1)
		}

		policies, err := autoscaling.GetPolicies(app, env)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		g, err := marathon.GroupDefinition(app, env)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		se := &scaleEvent{App: app, Env: env}
		for p, size := range rules {
			for _, a := range g.Apps {
				it := a.InstanceType()
				if it == p {
					for _, wp := range policies.QueueLengthScaled {
						if wp.Enabled && wp.AppName == app && wp.Environment == env && wp.AppType == it {
							fmt.Printf("autoscaling policy %s already enabled for %s-%s/%s\n", green(wp.Name), app, env, it)
							fmt.Printf("see: duncan autoscale list --app %s --env %s\n", app, env)
							os.Exit(1)
						}
					}
					for _, cp := range policies.CPUScaled {
						if cp.Enabled && cp.AppName == app && cp.Environment == env && cp.AppType == a.InstanceType() {
							fmt.Printf("autoscaling policy %s already enabled for %s-%s/%s\n", green(cp.Name), app, env, it)
							fmt.Printf("see: duncan autoscale list --app %s --env %s\n", app, env)
							os.Exit(1)
						}
					}
					se.Procs = append(se.Procs, &proc{
						InstanceType: p,
						Previous:     a.Instances,
						Current:      size,
					})
				}
			}
		}

		if promptScale(se) {
			id, err := marathon.Scale(g, rules)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if err := deployment.Watch(id); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			var msg string
			for _, proc := range se.Procs {
				msg = fmt.Sprintf(":whale: :scales: `%s` scaled from `%v` => `%v` instances\n", proc.InstanceType, proc.Previous, proc.Current)
			}
			if err := notify.Slack(viper.GetString("slack_webhook_url"), fmt.Sprintf("%s %s", app, env), msg); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(scaleCmd)

	scaleCmd.Flags().StringVarP(&app, "app", "a", "", "app to scale")
	scaleCmd.Flags().StringVarP(&env, "env", "e", "", "environment (stage, production)")
	scaleCmd.Flags().BoolVarP(&force, "force", "f", false, "bypass prompt before scaling")
}

func parseScaleRules(args []string) (map[string]int, error) {
	rules := map[string]int{}
	if len(args) == 0 {
		return rules, fmt.Errorf("wrong number of arguments")
	}

	for _, p := range args {
		s := strings.Split(p, "=")
		if len(s) != 2 {
			return rules, fmt.Errorf("%s not in proper format", p)
		}
		count, err := strconv.Atoi(s[1])
		if err != nil {
			return rules, fmt.Errorf("%s not in proper format", p)
		}
		if count < 1 {
			return rules, fmt.Errorf("scale count must be greater than 0")
		}
		rules[s[0]] = count
	}
	return rules, nil
}

func promptScale(se *scaleEvent) bool {
	if force {
		return true
	}
	white := color.New(color.FgWhite, color.Bold).SprintFunc()
	red := color.New(color.FgRed, color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
	fmt.Printf("You are about to scale:\n\n")
	fmt.Printf(white("  app: %s\n"), yellow(app))
	if env == "production" {
		fmt.Printf(white("  env: %s\n"), red(env))
	} else {
		fmt.Printf(white("  env: %s\n"), green(env))
	}
	for _, s := range se.Procs {
		fmt.Printf("  %s: %s => %s instances\n", white(s.InstanceType), cyan(s.Previous), yellow(s.Current))
	}
	fmt.Printf(`
NOTE: Manually scaling is often necessary for troubleshooting but consider instead using
an autoscaling policy once the conditions that trigger scaling and mininum/maximum
number of instances needed are well understood. The autoscaler is great at doing repetitive
tasks and won't forget to scale down when instances are no longer needed online.
`)

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
