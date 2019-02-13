package k8s

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"

	apiv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	yellow = color.New(color.FgYellow, color.Bold).SprintFunc()
	cyan   = color.New(color.FgCyan, color.Bold).SprintFunc()
	white  = color.New(color.FgWhite, color.Bold).SprintFunc()
	green  = color.New(color.FgGreen, color.Bold).SprintFunc()
)

type deploymentGroups map[string][][]string

// List displays k8s pods matching given app/env
func (k *KubeAPI) List(app, env string) error {
	if env == "" {
		env = "stage|production"
	}
	deploymentsClient := k.Client.AppsV1().Deployments(k.Namespace)
	ssClient := k.Client.AppsV1().StatefulSets(k.Namespace)

	deploymentList, err := deploymentsClient.List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	statefulSetList, err := ssClient.List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	groups := deploymentGroups{}
	groups = collectDeploymentGroups(deploymentList, app, env, groups)
	groups = collectStatefulSetGroups(statefulSetList, app, env, groups)

	for k, v := range groups {
		fmt.Println(green(k))
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Tag", "Instances", "CPU", "Mem"})
		table.AppendBulk(v)
		table.Render()
	}

	return nil
}

func collectDeploymentGroups(deploymentList *apiv1.DeploymentList, app string, env string, groups deploymentGroups) deploymentGroups {
	for _, item := range deploymentList.Items {
		group := item.Spec.Template.ObjectMeta.Labels["group"]
		groupEnv := item.Spec.Template.ObjectMeta.Labels["env"]
		replicas := item.Status.Replicas
		for _, container := range item.Spec.Template.Spec.Containers {
			groups = addContainerToGroup(groups, app, env, group, groupEnv, replicas, container)
		}
	}

	return groups
}

func collectStatefulSetGroups(deploymentList *apiv1.StatefulSetList, app string, env string, groups deploymentGroups) deploymentGroups {
	for _, item := range deploymentList.Items {
		group := item.Spec.Template.ObjectMeta.Labels["group"]
		groupEnv := item.Spec.Template.ObjectMeta.Labels["env"]
		replicas := item.Status.Replicas
		for _, container := range item.Spec.Template.Spec.Containers {
			groups = addContainerToGroup(groups, app, env, group, groupEnv, replicas, container)
		}
	}

	return groups
}

func addContainerToGroup(groups deploymentGroups, app, env, group, groupEnv string, replicas int32, container corev1.Container) deploymentGroups {
	var data = make([][]string, 10)
	parts := strings.Split(container.Image, ":")
	var tag string
	if len(parts) > 1 {
		tag = parts[1]
	}
	cpu := container.Resources.Limits["cpu"]
	mem := container.Resources.Limits["memory"]
	data = append(data, []string{
		cyan(container.Name),
		white(tag),
		yellow(replicas),
		cyan(cpu.String()),
		cyan(mem.String()),
	})

	parts = strings.Split(group, "-")
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

	return groups
}
