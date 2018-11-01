package k8s

import (
	"fmt"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

// CurrentTag fetches the currently deployed docker image tag for
// given app and env if it exists. First checks Kubernetes Deployment API
// and then Stateful Sets API
func (k *KubeAPI) CurrentTag(app, env, repo string) (string, error) {
	deploymentsClient := k.Client.AppsV1().Deployments(k.Namespace)
	deploymentList, err := deploymentsClient.List(metav1.ListOptions{})
	if err != nil {
		return "", err
	}
	for _, item := range deploymentList.Items {
		tag := findTag(app, env, item.Spec.Template)
		if tag != "" {
			return tag, nil
		}
	}

	ssClient := k.Client.AppsV1().StatefulSets(k.Namespace)
	ssList, err := ssClient.List(metav1.ListOptions{})
	if err != nil {
		return "", err
	}
	for _, item := range ssList.Items {
		tag := findTag(app, env, item.Spec.Template)
		if tag != "" {
			return tag, nil
		}
	}

	return "", fmt.Errorf("could not find tag for %s-%s", app, env)
}

// Deploy updates docker image tag for a given k8s deployment
func (k *KubeAPI) Deploy(app, env, tag, repo string) error {
	if err := k.updateDeployment(app, env, tag, repo); err != nil {
		return err
	}

	return k.updateStatefulSet(app, env, tag, repo)
}

func findTag(app, env string, template corev1.PodTemplateSpec) string {
	group := template.ObjectMeta.Labels["group"]
	containers := template.Spec.Containers

	if group == fmt.Sprintf("%s-%s", app, env) {
		for _, container := range containers {
			parts := strings.Split(container.Image, ":")
			if len(parts) != 2 {
				continue
			}
			return parts[1]
		}
	}

	return ""
}

func (k *KubeAPI) updateDeployment(app, env, tag, repo string) error {
	deploymentsClient := k.Client.AppsV1().Deployments(k.Namespace)

	list, err := deploymentsClient.List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	var toUpdate []appsv1.Deployment
	for _, item := range list.Items {
		group := item.Spec.Template.ObjectMeta.Labels["group"]
		if group == fmt.Sprintf("%s-%s", app, env) {
			toUpdate = append(toUpdate, item)
		}
	}

	for _, deployment := range toUpdate {
		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			item, err := deploymentsClient.Get(deployment.ObjectMeta.Name, metav1.GetOptions{})
			if err != nil {
				return err
			}
			image := item.Spec.Template.Spec.Containers[0].Image
			parts := strings.Split(image, ":")
			item.Spec.Template.Spec.Containers[0].Image = fmt.Sprintf("%s:%s", parts[0], tag)
			_, err = deploymentsClient.Update(item)
			return err
		})
		if retryErr != nil {
			return fmt.Errorf("deploy failed: %v", retryErr)
		}
	}

	return nil
}

func (k *KubeAPI) updateStatefulSet(app, env, tag, repo string) error {
	ssClient := k.Client.AppsV1().StatefulSets(k.Namespace)

	ssList, err := ssClient.List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	var toUpdate []appsv1.StatefulSet
	for _, item := range ssList.Items {
		group := item.Spec.Template.ObjectMeta.Labels["group"]
		if group == fmt.Sprintf("%s-%s", app, env) {
			toUpdate = append(toUpdate, item)
		}
	}

	for _, deployment := range toUpdate {
		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			item, err := ssClient.Get(deployment.ObjectMeta.Name, metav1.GetOptions{})
			if err != nil {
				return err
			}
			image := item.Spec.Template.Spec.Containers[0].Image
			parts := strings.Split(image, ":")
			item.Spec.Template.Spec.Containers[0].Image = fmt.Sprintf("%s:%s", parts[0], tag)
			_, err = ssClient.Update(item)
			return err
		})
		if retryErr != nil {
			return fmt.Errorf("deploy failed: %v", retryErr)
		}
	}

	return nil
}
