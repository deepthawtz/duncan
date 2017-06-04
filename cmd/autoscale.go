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
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/betterdoctor/slythe/policy"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
			policies, err := getPolicies(app, env)
			if err != nil {
				fmt.Printf("failed to fetch policies: %v\n", err)
				os.Exit(1)
			}

			fmt.Println(green("CPU Scaling Policies"))
			for _, cp := range policies.CPUScaled {
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
				fmt.Fprintln(w, white("Policy Name \t"), cyan(cp.Name))
				fmt.Fprintln(w, white("App \t"), green(cp.AppName))
				fmt.Fprintln(w, white("App Type \t"), green(cp.AppType))
				fmt.Fprintln(w, white("Env \t"), green(cp.Environment))
				fmt.Fprintln(w, white("Min Instances \t"), white(cp.MinInstances))
				fmt.Fprintln(w, white("Max Instances \t"), white(cp.MaxInstances))
				fmt.Fprintln(w, white("Scale Up By \t"), white(cp.ScaleUpBy))
				fmt.Fprintln(w, white("Scale Up By \t"), white(cp.ScaleDownBy))
				fmt.Fprintln(w, white("Up Threshold \t"), yellow(fmt.Sprintf("%d%%", cp.UpThreshold)))
				fmt.Fprintln(w, white("Down Threshold \t"), yellow(fmt.Sprintf("%d%%", cp.DownThreshold)))
				w.Flush()
				fmt.Println("-------------------------------------")
			}

			fmt.Println()
			fmt.Println(green("Worker Scaling Policies"))
			for _, wp := range policies.QueueLengthScaled {
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
				fmt.Fprintln(w, white("Policy Name \t"), cyan(wp.Name))
				fmt.Fprintln(w, white("App \t"), green(wp.AppName))
				fmt.Fprintln(w, white("App Type \t"), green(wp.AppType))
				fmt.Fprintln(w, white("Env \t"), green(wp.Environment))
				fmt.Fprintln(w, white("Min Instances \t"), white(wp.MinInstances))
				fmt.Fprintln(w, white("Max Instances \t"), white(wp.MaxInstances))
				fmt.Fprintln(w, white("Scale Up By \t"), white(wp.ScaleUpBy))
				fmt.Fprintln(w, white("Scale Down By \t"), white(wp.ScaleDownBy))
				fmt.Fprintln(w, white("Up Threshold \t"), yellow(wp.UpThreshold))
				fmt.Fprintln(w, white("Down Threshold \t"), yellow(wp.DownThreshold))
				fmt.Fprintln(w, white("Redis URL \t"), green(wp.RedisURL))
				fmt.Fprintln(w, white("Queues \t"), cyan(wp.Queues))
				w.Flush()
				fmt.Println("-------------------------------------")
			}
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
			fmt.Println("implement me")
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
			fmt.Println("implement me")
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
			fmt.Println("implement me")
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
			fmt.Println("implement me")
		},
	}
)

func init() {
	autoscaleCPUCmd.PersistentFlags().StringVarP(&policyName, "policy-name", "", "", "name for autoscaling policy")
	autoscaleCPUCmd.PersistentFlags().StringVarP(&app, "app", "a", "", "app to manage autoscaling policy for")
	autoscaleCPUCmd.PersistentFlags().StringVarP(&env, "env", "e", "", "app environment to manage autoscaling policy for")
	autoscaleCPUCmd.PersistentFlags().StringVarP(&appType, "app-type", "", "", "app type for autoscaling policy")
	autoscaleCPUCmd.PersistentFlags().IntVarP(&min, "min-instances", "", 1, "minimum number of instances for autoscaling policy")
	autoscaleCPUCmd.PersistentFlags().IntVarP(&max, "max-instances", "", 1, "maximum number of instances for autoscaling policy")
	autoscaleCPUCmd.PersistentFlags().IntVarP(&upBy, "scale-up-by", "", 1, "number of instances to scale up by when up threshold exceeded")
	autoscaleCPUCmd.PersistentFlags().IntVarP(&downBy, "scale-down-by", "", 1, "number of instances to scale down by when below down threshold value")
	autoscaleCPUCmd.PersistentFlags().IntVarP(&upThreshold, "up-threshold", "", 0, "aggregate CPU percent to scale up on")
	autoscaleCPUCmd.PersistentFlags().IntVarP(&downThreshold, "down-threshold", "", 0, "aggregate CPU percent to scale down on")
	autoscaleCPUCmd.PersistentFlags().IntVarP(&checkFreqSecs, "check-frequency-secs", "", 30, "frequency the autoscaling policy will be repeatedly checked (minimum: 10 seconds)")

	autoscaleWorkerCmd.PersistentFlags().StringVarP(&policyName, "policy-name", "", "", "name for autoscaling policy")
	autoscaleWorkerCmd.PersistentFlags().StringVarP(&app, "app", "a", "", "app to manage autoscaling policy for")
	autoscaleWorkerCmd.PersistentFlags().StringVarP(&env, "env", "e", "", "app environment to manage autoscaling policy for")
	autoscaleWorkerCmd.PersistentFlags().StringVarP(&appType, "app-type", "", "", "app type for autoscaling policy")
	autoscaleWorkerCmd.PersistentFlags().StringVarP(&redisURL, "redis-url", "", "", "redis URL for autoscaling policy")
	autoscaleWorkerCmd.PersistentFlags().StringVarP(&queues, "queues", "", "", "redis queues for worker autoscaling policy")
	autoscaleWorkerCmd.PersistentFlags().IntVarP(&min, "min-instances", "", 1, "minimum number of instances for autoscaling policy")
	autoscaleWorkerCmd.PersistentFlags().IntVarP(&max, "max-instances", "", 1, "maximum number of instances for autoscaling policy")
	autoscaleWorkerCmd.PersistentFlags().IntVarP(&upBy, "scale-up-by", "", 1, "number of instances to scale up by when up threshold exceeded")
	autoscaleWorkerCmd.PersistentFlags().IntVarP(&downBy, "scale-down-by", "", 1, "number of instances to scale down by when below down threshold value")
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

func getPolicies(app, env string) (*policy.Policies, error) {
	policies := &policy.Policies{}
	resp, err := http.Get(viper.GetString("SLYTHE_HOST") + "/")
	if err != nil {
		return policies, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(policies); err != nil {
		return policies, err
	}
	if app == "" && env == "" {
		return policies, nil
	}
	if env == "" {
		env = "stage|production"
	}
	fp := &policy.Policies{}
	for _, e := range strings.Split(env, "|") {
		for _, cp := range policies.CPUScaled {
			if app == "" && cp.Environment == e {
				fp.CPUScaled = append(fp.CPUScaled, cp)
			}
			if app != "" && cp.AppName == app && cp.Environment == e {
				fp.CPUScaled = append(fp.CPUScaled, cp)
			}
		}
	}
	for _, e := range strings.Split(env, "|") {
		for _, cp := range policies.QueueLengthScaled {
			if app == "" && cp.Environment == e {
				fp.QueueLengthScaled = append(fp.QueueLengthScaled, cp)
			}
			if app != "" && cp.AppName == app && cp.Environment == e {
				fp.QueueLengthScaled = append(fp.QueueLengthScaled, cp)
			}
		}
	}
	return fp, nil
}
