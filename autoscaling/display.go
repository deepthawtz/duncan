package autoscaling

import (
	"fmt"
	"os"
	"text/tabwriter"

	pb "github.com/betterdoctor/slythe/rpc"
	"github.com/fatih/color"
)

var (
	yellow = color.New(color.FgYellow, color.Bold).SprintFunc()
	cyan   = color.New(color.FgCyan, color.Bold).SprintFunc()
	white  = color.New(color.FgWhite, color.Bold).SprintFunc()
	green  = color.New(color.FgGreen, color.Bold).SprintFunc()
	red    = color.New(color.FgRed, color.Bold).SprintFunc()
)

// DisplayCPUPolicies prints the CPU policies
func DisplayCPUPolicies(policies *pb.Policies) {
	if len(policies.CpuPolicies) > 0 {
		fmt.Printf(green("CPU Scaling Policies\n\n"))
		for _, cp := range policies.CpuPolicies {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
			fmt.Fprintln(w, white("Policy Name \t"), cyan(cp.Name))
			fmt.Fprintln(w, white("App \t"), green(cp.AppName))
			fmt.Fprintln(w, white("App Type \t"), green(cp.AppType))
			fmt.Fprintln(w, white("Env \t"), green(cp.Environment))
			fmt.Fprintln(w, white("Min Instances \t"), white(cp.MinInstances))
			fmt.Fprintln(w, white("Max Instances \t"), white(cp.MaxInstances))
			fmt.Fprintln(w, white("Scale Up By \t"), white(cp.ScaleUpBy))
			fmt.Fprintln(w, white("Scale Down By \t"), white(cp.ScaleDownBy))
			fmt.Fprintln(w, white("Up Threshold \t"), yellow(fmt.Sprintf("%d%%", cp.UpThreshold)))
			fmt.Fprintln(w, white("Down Threshold \t"), yellow(fmt.Sprintf("%d%%", cp.DownThreshold)))
			fmt.Fprintln(w, white("Check Frequency Secs \t"), white(fmt.Sprintf("%d", cp.CheckFrequencySecs)))
			if cp.Enabled {
				fmt.Fprintln(w, white("Enabled \t"), green("true"))
			} else {
				fmt.Fprintln(w, white("Enabled \t"), red("false"))
			}
			w.Flush()
			fmt.Println("-------------------------------------")
		}
	}
}

// DisplayWorkerPolicies prints the Worker policies
func DisplayWorkerPolicies(policies *pb.Policies) {
	if len(policies.WorkerPolicies) > 0 {
		fmt.Printf(green("Worker Scaling Policies\n\n"))
		for _, wp := range policies.WorkerPolicies {
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
			fmt.Fprintln(w, white("Redis URL \t"), green(wp.RedisUrl))
			fmt.Fprintln(w, white("Queues \t"), cyan(wp.Queues))
			fmt.Fprintln(w, white("Check Frequency Secs \t"), white(fmt.Sprintf("%d", wp.CheckFrequencySecs)))
			if wp.Enabled {
				fmt.Fprintln(w, white("Enabled \t"), green("true"))
			} else {
				fmt.Fprintln(w, white("Enabled \t"), red("false"))
			}
			w.Flush()
			fmt.Println("-------------------------------------")
		}
	}
}
