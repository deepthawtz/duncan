package k8s

import (
	"fmt"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

// CurrentTag fetches the currently deployed docker image tag for
// given app and env
func CurrentTag(app, env, repo string) (string, error) {
	clientset, err := newClient()

	deploymentsClient := clientset.AppsV1().Deployments("pipeline")

	item, err := deploymentsClient.Get(fmt.Sprintf("%s-%s-web", app, env), metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	for _, container := range item.Spec.Template.Spec.Containers {
		parts := strings.Split(container.Image, ":")
		if len(parts) != 2 {
			return "", fmt.Errorf("could not find tag for %s-%s", app, env)
		}
		tag := parts[1]
		return tag, nil
	}

	return "", fmt.Errorf("could not find tag for %s-%s", app, env)
}

// Deploy updates docker image tag for a given k8s deployment
func Deploy(app, env, tag, repo string) error {
	clientset, err := newClient()
	if err != nil {
		return err
	}

	deploymentsClient := clientset.AppsV1().Deployments("pipeline")

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
