package k8s

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// List displays k8s pods matching given app/env
func List(app, env string) error {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	if env == "" {
		env = "stage|production"
	}
	deploymentsClient := clientset.AppsV1().Deployments("pipeline")

	list, err := deploymentsClient.List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	white := color.New(color.FgWhite, color.Bold).SprintFunc()
	green := color.New(color.FgGreen, color.Bold).SprintFunc()

	groups := map[string][][]string{}
	for _, item := range list.Items {
		group := item.Spec.Template.ObjectMeta.Labels["group"]
		groupEnv := item.Spec.Template.ObjectMeta.Labels["env"]
		replicas := item.Status.Replicas
		for _, container := range item.Spec.Template.Spec.Containers {
			var data = make([][]string, 10)
			data = append(data, []string{
				cyan(container.Name),
				white(strings.Split(container.Image, ":")[1]),
				yellow(replicas),
			})

			parts := strings.Split(group, "-")
			a := strings.Join(parts[:len(parts)-1], "-")

			for _, e := range strings.Split(env, "|") {
				if groupEnv == e {
					if app == "" || a == app {
						for _, d := range data {
							groups[group] = append(groups[group], d)
						}
					}
				}
			}
		}
	}

	for k, v := range groups {
		fmt.Println(green(k))
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Tag", "Instances"})
		table.AppendBulk(v)
		table.Render()
	}

	return nil
}
