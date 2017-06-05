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
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/betterdoctor/duncan/autoscaling"
	"github.com/betterdoctor/slythe/policy"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	yellow = color.New(color.FgYellow, color.Bold).SprintFunc()
	cyan   = color.New(color.FgCyan, color.Bold).SprintFunc()
	white  = color.New(color.FgWhite, color.Bold).SprintFunc()
	green  = color.New(color.FgGreen, color.Bold).SprintFunc()

	// args for autoscaling policy commands
	policyName, appType, redisURL, queues                             string
	min, max, upBy, downBy, checkFreqSecs, upThreshold, downThreshold int

	autoscaleCmd = &cobra.Command{
		Use:   "autoscale",
		Short: "Commands to manage autoscaling policies",
		Long:  "Commands to create and modify autoscaling policies including enable/disable.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("must provide autoscale subcommand, see: duncan autoscale -h")
			os.Exit(1)
		},
	}

	autoscaleListCmd = &cobra.Command{
		Use:   "list",
		Short: "List autoscaling policies",
		Run: func(cmd *cobra.Command, args []string) {
			policies, err := autoscaling.GetPolicies(app, env)
			if err != nil {
				fmt.Printf("failed to fetch policies: %v\n", err)
				os.Exit(1)
			}

			autoscaling.DisplayPolicies(policies)
		},
	}

	autoscaleCPUCmd = &cobra.Command{
		Use:   "cpu",
		Short: "Commands to manage CPU-based autoscaling policies",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("must provide autoscale subcommand, see: duncan autoscale cpu -h")
			os.Exit(1)
		},
	}

	autoscaleWorkerCmd = &cobra.Command{
		Use:   "worker",
		Short: "Commands to manage worker autoscaling policies",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("must provide autoscale subcommand, see: duncan autoscale worker -h")
			os.Exit(1)
		},
	}

	autoscaleWorkerCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create new worker autoscaling policy",
		Long: `Create new Worker autoscaling policy.

Examples:

$ duncan autoscale worker create --app myapp --env production --policy-name MyAppProductionWorker \
  --app-type worker --min-instances 1 --max-instances 20 --up-threshold 5000 --down-threshold 1000 \
  --redis-url redis://yo.dawg:6379/2 --queues queue:important_stuff,queue:cat_photos,queue:junk \
  --scale-up-by 3 --scale-down-by 1 --check-frequency-secs 30
		`,
		Run: func(cmd *cobra.Command, args []string) {
			validateCreatePolicyFlags("worker")
			if promptCreateScalingPolicy() {
				gp := newGenericPolicy()
				wp := &policy.Worker{
					GenericPolicy: gp,
					RedisURL:      redisURL,
					Queues:        queues,
				}
				if err := autoscaling.CreateWorkerPolicy(wp); err != nil {
					fmt.Printf("failed to create autoscaling policy: %v\n", err)
					os.Exit(1)
				}
				fmt.Printf("autoscaling policy %s created\n", green(wp.Name))
			}
		},
	}

	autoscaleCPUCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create new CPU autoscaling policy",
		Long: `Create new CPU autoscaling policy.

Examples:

$ duncan autoscale cpu create --app myapp --env production --policy-name MyAppProductionWeb \
  --app-type web --min-instances 1 --max-instances 10 --scale-up-by 2 --scale-down-by 1 \
  --check-frequency-secs 30 --up-threshold 50 --down-threshold 5
		`,
		Run: func(cmd *cobra.Command, args []string) {
			validateCreatePolicyFlags("cpu")
			if promptCreateScalingPolicy() {
				gp := newGenericPolicy()
				cp := &policy.CPU{
					GenericPolicy: gp,
				}
				if err := autoscaling.CreateCPUPolicy(cp); err != nil {
					fmt.Printf("failed to create autoscaling policy: %v\n", err)
					os.Exit(1)
				}
				fmt.Printf("autoscaling policy %s created\n", green(cp.Name))
			}
		},
	}

	autoscaleWorkerUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "update worker autoscaling policy",
		Long: `Update Worker autoscaling policy.

Examples:

$ duncan autoscale worker update --policy-name MyAppProductionWorker --max-instances 200
		`,
		Run: func(cmd *cobra.Command, args []string) {
			if policyName == "" {
				fmt.Println("update command requires --policy-name flag")
				os.Exit(1)
			}
			policies, err := autoscaling.GetPolicies(app, env)
			if err != nil {
				fmt.Printf("failed to fetch policies: %s\n", err)
				os.Exit(1)
			}
			wp := &policy.Worker{}
			for _, w := range policies.QueueLengthScaled {
				if w.Name == policyName {
					wp = w
				}
			}
			if wp.Name == "" {
				fmt.Printf("could not find policy %s to update\n", green(policyName))
				os.Exit(1)
			}
			for k, v := range workerStringFlags() {
				if k == "--app" && v != "" {
					wp.AppName = v
				}
				if k == "--env" && v != "" {
					wp.Environment = v
				}
				if k == "--app-type" && v != "" {
					wp.AppType = v
				}
				if k == "--redis-url" && v != "" {
					wp.RedisURL = v
				}
				if k == "--queues" && v != "" {
					wp.Queues = v
				}
			}
			for k, v := range policyIntFlags() {
				if k == "--min-instances" && v != 0 {
					wp.MinInstances = v
				}
				if k == "--max-instances" && v != 0 {
					wp.MaxInstances = v
				}
				if k == "--scale-up-by" && v != 0 {
					wp.ScaleUpBy = v
				}
				if k == "--scale-down-by" && v != 0 {
					wp.ScaleDownBy = v
				}
				if k == "--up-threshold" && v != 0 {
					wp.UpThreshold = v
				}
				if k == "--down-threshold" && v != 0 {
					wp.DownThreshold = v
				}
				if k == "--check-frequency-secs" && v != 0 {
					wp.CheckFrequencySecs = v
				}
			}
			if promptUpdateWorkerScalingPolicy(wp) {
				if err := autoscaling.UpdateWorkerPolicy(wp); err != nil {
					fmt.Printf("failed to update autoscaling policy: %v\n", err)
					os.Exit(1)
				}
			}
			fmt.Printf("autoscaling policy %s updated\n", green(wp.Name))
		},
	}

	autoscaleCPUUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "update CPU autoscaling policy",
		Long: `Update CPU autoscaling policy.

Examples:

$ duncan autoscale cpu update --policy-name MyAppProductionWeb --max-instances 50
		`,
		Run: func(cmd *cobra.Command, args []string) {
			if policyName == "" {
				fmt.Println("update command requires --policy-name flag")
				os.Exit(1)
			}
			policies, err := autoscaling.GetPolicies(app, env)
			if err != nil {
				fmt.Printf("failed to fetch policies: %s\n", err)
				os.Exit(1)
			}
			cp := &policy.CPU{}
			for _, c := range policies.CPUScaled {
				if c.Name == policyName {
					cp = c
				}
			}
			if cp.Name == "" {
				fmt.Printf("could not find policy %s to update\n", green(policyName))
				os.Exit(1)
			}
			for k, v := range cpuStringFlags() {
				if k == "--app" && v != "" {
					cp.AppName = v
				}
				if k == "--env" && v != "" {
					cp.Environment = v
				}
				if k == "--app-type" && v != "" {
					cp.AppType = v
				}
			}
			for k, v := range policyIntFlags() {
				if k == "--min-instances" && v != 0 {
					cp.MinInstances = v
				}
				if k == "--max-instances" && v != 0 {
					cp.MaxInstances = v
				}
				if k == "--scale-up-by" && v != 0 {
					cp.ScaleUpBy = v
				}
				if k == "--scale-down-by" && v != 0 {
					cp.ScaleDownBy = v
				}
				if k == "--up-threshold" && v != 0 {
					cp.UpThreshold = v
				}
				if k == "--down-threshold" && v != 0 {
					cp.DownThreshold = v
				}
				if k == "--check-frequency-secs" && v != 0 {
					cp.CheckFrequencySecs = v
				}
			}
			if promptUpdateCPUScalingPolicy(cp) {
				if err := autoscaling.UpdateCPUPolicy(cp); err != nil {
					fmt.Printf("failed to update autoscaling policy: %v\n", err)
					os.Exit(1)
				}
			}
			fmt.Printf("autoscaling policy %s updated\n", green(cp.Name))
		},
	}
)

func init() {
	autoscaleCPUCmd.PersistentFlags().StringVarP(&policyName, "policy-name", "", "", "name for autoscaling policy")
	autoscaleCPUCmd.PersistentFlags().StringVarP(&app, "app", "a", "", "app to manage autoscaling policy for")
	autoscaleCPUCmd.PersistentFlags().StringVarP(&env, "env", "e", "", "app environment to manage autoscaling policy for")
	autoscaleCPUCmd.PersistentFlags().StringVarP(&appType, "app-type", "", "", "app type for autoscaling policy")
	autoscaleCPUCmd.PersistentFlags().IntVarP(&min, "min-instances", "", 0, "minimum number of instances for autoscaling policy")
	autoscaleCPUCmd.PersistentFlags().IntVarP(&max, "max-instances", "", 0, "maximum number of instances for autoscaling policy")
	autoscaleCPUCmd.PersistentFlags().IntVarP(&upBy, "scale-up-by", "", 0, "number of instances to scale up by when up threshold exceeded")
	autoscaleCPUCmd.PersistentFlags().IntVarP(&downBy, "scale-down-by", "", 0, "number of instances to scale down by when below down threshold value")
	autoscaleCPUCmd.PersistentFlags().IntVarP(&upThreshold, "up-threshold", "", 0, "aggregate CPU percent to scale up on")
	autoscaleCPUCmd.PersistentFlags().IntVarP(&downThreshold, "down-threshold", "", 0, "aggregate CPU percent to scale down on")
	autoscaleCPUCmd.PersistentFlags().IntVarP(&checkFreqSecs, "check-frequency-secs", "", 30, "frequency the autoscaling policy will be repeatedly checked (minimum: 10 seconds)")

	autoscaleWorkerCmd.PersistentFlags().StringVarP(&policyName, "policy-name", "", "", "name for autoscaling policy")
	autoscaleWorkerCmd.PersistentFlags().StringVarP(&app, "app", "a", "", "app to manage autoscaling policy for")
	autoscaleWorkerCmd.PersistentFlags().StringVarP(&env, "env", "e", "", "app environment to manage autoscaling policy for")
	autoscaleWorkerCmd.PersistentFlags().StringVarP(&appType, "app-type", "", "", "app type for autoscaling policy")
	autoscaleWorkerCmd.PersistentFlags().StringVarP(&redisURL, "redis-url", "", "", "redis URL for autoscaling policy")
	autoscaleWorkerCmd.PersistentFlags().StringVarP(&queues, "queues", "", "", "redis queues for worker autoscaling policy")
	autoscaleWorkerCmd.PersistentFlags().IntVarP(&min, "min-instances", "", 0, "minimum number of instances for autoscaling policy")
	autoscaleWorkerCmd.PersistentFlags().IntVarP(&max, "max-instances", "", 0, "maximum number of instances for autoscaling policy")
	autoscaleWorkerCmd.PersistentFlags().IntVarP(&upBy, "scale-up-by", "", 0, "number of instances to scale up by when up threshold exceeded")
	autoscaleWorkerCmd.PersistentFlags().IntVarP(&downBy, "scale-down-by", "", 0, "number of instances to scale down by when below down threshold value")
	autoscaleWorkerCmd.PersistentFlags().IntVarP(&upThreshold, "up-threshold", "", 0, "redis queue size to scale up on")
	autoscaleWorkerCmd.PersistentFlags().IntVarP(&downThreshold, "down-threshold", "", 0, "redis queue size to scale down on")
	autoscaleWorkerCmd.PersistentFlags().IntVarP(&checkFreqSecs, "check-frequency-secs", "", 30, "frequency the autoscaling policy will be repeatedly checked (minimum: 10 seconds)")

	RootCmd.AddCommand(autoscaleCmd)
	autoscaleCmd.AddCommand(autoscaleListCmd)
	autoscaleListCmd.Flags().StringVarP(&app, "app", "a", "", "app to list autoscaling policies for")
	autoscaleListCmd.Flags().StringVarP(&env, "env", "e", "", "env to list autoscaling policies for")
	autoscaleCmd.AddCommand(autoscaleWorkerCmd)
	autoscaleCmd.AddCommand(autoscaleCPUCmd)
	autoscaleWorkerCmd.AddCommand(autoscaleWorkerCreateCmd)
	autoscaleWorkerCmd.AddCommand(autoscaleWorkerUpdateCmd)
	autoscaleCPUCmd.AddCommand(autoscaleCPUCreateCmd)
	autoscaleCPUCmd.AddCommand(autoscaleCPUUpdateCmd)
}

func cpuStringFlags() map[string]string {
	return map[string]string{
		"--policy-name": policyName,
		"--app":         app,
		"--env":         env,
		"--app-type":    appType,
	}
}

func workerStringFlags() map[string]string {
	m := cpuStringFlags()
	m["--redis-url"] = redisURL
	m["--queues"] = queues
	return m
}

func policyIntFlags() map[string]int {
	return map[string]int{
		"--min-instances":        min,
		"--max-instances":        max,
		"--scale-up-by":          upBy,
		"--scale-down-by":        downBy,
		"--up-threshold":         upThreshold,
		"--down-threshold":       downThreshold,
		"--check-frequency-secs": checkFreqSecs,
	}
}

func validateCreatePolicyFlags(policyType string) {
	var stringFlagsMissing, intFlagsMissing bool
	var rs map[string]string
	if policyType == "worker" {
		rs = workerStringFlags()
	} else {
		rs = cpuStringFlags()
	}
	for k, v := range rs {
		if v == "" {
			stringFlagsMissing = true
			fmt.Printf("%s flag is required\n", k)
		}
	}
	ri := policyIntFlags()
	for k, v := range ri {
		if v == 0 {
			intFlagsMissing = true
			fmt.Printf("%s flag is required\n", k)
		}
	}
	if stringFlagsMissing || intFlagsMissing {
		fmt.Println("missing required flags")
		os.Exit(1)
	}
}

func promptCreateScalingPolicy() bool {
	fmt.Println("You are about to create the following autoscaling policy")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, white("Policy Name \t"), cyan(policyName))
	fmt.Fprintln(w, white("App \t"), green(app))
	fmt.Fprintln(w, white("App Type \t"), green(appType))
	fmt.Fprintln(w, white("Env \t"), green(env))
	fmt.Fprintln(w, white("Min Instances \t"), white(min))
	fmt.Fprintln(w, white("Max Instances \t"), white(max))
	fmt.Fprintln(w, white("Scale Up By \t"), white(upBy))
	fmt.Fprintln(w, white("Scale Down By \t"), white(downBy))
	fmt.Fprintln(w, white("Up Threshold \t"), yellow(upThreshold))
	fmt.Fprintln(w, white("Down Threshold \t"), yellow(downThreshold))
	if redisURL != "" && queues != "" {
		fmt.Fprintln(w, white("Redis URL \t"), green(redisURL))
		fmt.Fprintln(w, white("Queues \t"), cyan(queues))
	}
	w.Flush()
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(white("\nare you sure? (yes/no): "))
	resp, _ := reader.ReadString('\n')

	resp = strings.TrimSpace(resp)
	if resp != "yes" {
		return false
	}
	return true
}

func promptUpdateWorkerScalingPolicy(wp *policy.Worker) bool {
	fmt.Println("You are about to update the following autoscaling policy")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, white("Policy Name \t"), cyan(wp.Name))
	fmt.Fprintln(w, white("App \t"), green(wp.AppName))
	fmt.Fprintln(w, white("App Type \t"), green(wp.AppType))
	fmt.Fprintln(w, white("Env \t"), green(wp.Environment))
	fmt.Fprintln(w, white("Redis URL \t"), green(wp.RedisURL))
	fmt.Fprintln(w, white("Queues \t"), cyan(wp.Queues))
	fmt.Fprintln(w, white("Min Instances \t"), white(wp.MinInstances))
	fmt.Fprintln(w, white("Max Instances \t"), white(wp.MaxInstances))
	fmt.Fprintln(w, white("Scale Up By \t"), white(wp.ScaleUpBy))
	fmt.Fprintln(w, white("Scale Down By \t"), white(wp.ScaleDownBy))
	fmt.Fprintln(w, white("Up Threshold \t"), yellow(wp.UpThreshold))
	fmt.Fprintln(w, white("Down Threshold \t"), yellow(wp.DownThreshold))
	fmt.Fprintln(w, white("Check Frequency Secs \t"), yellow(wp.CheckFrequencySecs))
	w.Flush()
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(white("\nare you sure? (yes/no): "))
	resp, _ := reader.ReadString('\n')

	resp = strings.TrimSpace(resp)
	if resp != "yes" {
		return false
	}
	return true
}

func promptUpdateCPUScalingPolicy(cp *policy.CPU) bool {
	fmt.Println("You are about to update the following autoscaling policy")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, white("Policy Name \t"), cyan(cp.Name))
	fmt.Fprintln(w, white("App \t"), green(cp.AppName))
	fmt.Fprintln(w, white("App Type \t"), green(cp.AppType))
	fmt.Fprintln(w, white("Env \t"), green(cp.Environment))
	fmt.Fprintln(w, white("Min Instances \t"), white(cp.MinInstances))
	fmt.Fprintln(w, white("Max Instances \t"), white(cp.MaxInstances))
	fmt.Fprintln(w, white("Scale Up By \t"), white(cp.ScaleUpBy))
	fmt.Fprintln(w, white("Scale Down By \t"), white(cp.ScaleDownBy))
	fmt.Fprintln(w, white("Up Threshold \t"), yellow(cp.UpThreshold))
	fmt.Fprintln(w, white("Down Threshold \t"), yellow(cp.DownThreshold))
	fmt.Fprintln(w, white("Check Frequency Secs \t"), yellow(cp.CheckFrequencySecs))
	w.Flush()
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(white("\nare you sure? (yes/no): "))
	resp, _ := reader.ReadString('\n')

	resp = strings.TrimSpace(resp)
	if resp != "yes" {
		return false
	}
	return true
}

func newGenericPolicy() policy.GenericPolicy {
	return policy.GenericPolicy{
		Name:               policyName,
		AppName:            app,
		AppType:            appType,
		Environment:        env,
		MinInstances:       min,
		MaxInstances:       max,
		UpThreshold:        upThreshold,
		DownThreshold:      downThreshold,
		ScaleUpBy:          upBy,
		ScaleDownBy:        downBy,
		CheckFrequencySecs: checkFreqSecs,
		Enabled:            true,
	}
}
