package autoscaling

import (
	"context"
	"fmt"
	"net/http"

	pb "github.com/betterdoctor/slythe/rpc"
	"github.com/spf13/viper"
)

// GetPolicies returns all autoscaling policies optionally filtering
// if app or env are not empty string
func GetPolicies(app, env string) (*pb.Policies, error) {
	client := pb.NewAutoscalerProtobufClient(viper.GetString("SLYTHE_HOST"), &http.Client{})
	policies, err := client.ListPolicies(context.Background(), &pb.Filter{App: app, Env: env})
	if err != nil {
		return &pb.Policies{}, fmt.Errorf("failed to fetch policies: %v", err)
	}
	return policies, nil
}

// CreateWorkerPolicy creates an autoscaling worker policy
func CreateWorkerPolicy(wp *pb.WorkerPolicy) error {
	client := pb.NewAutoscalerProtobufClient(viper.GetString("SLYTHE_HOST"), &http.Client{})
	_, err := client.CreateWorkerAutoscalingPolicy(context.Background(), wp)
	return err
}

// UpdateWorkerPolicy updates an autoscaling worker policy
func UpdateWorkerPolicy(wp *pb.WorkerPolicy) error {
	client := pb.NewAutoscalerProtobufClient(viper.GetString("SLYTHE_HOST"), &http.Client{})
	_, err := client.UpdateWorkerAutoscalingPolicy(context.Background(), wp)
	return err
}

// CreateCPUPolicy creates an autoscaling worker policy
func CreateCPUPolicy(cp *pb.CPUPolicy) error {
	client := pb.NewAutoscalerProtobufClient(viper.GetString("SLYTHE_HOST"), &http.Client{})
	_, err := client.CreateCPUAutoscalingPolicy(context.Background(), cp)
	return err
}

// UpdateCPUPolicy creates an autoscaling worker policy
func UpdateCPUPolicy(cp *pb.CPUPolicy) error {
	client := pb.NewAutoscalerProtobufClient(viper.GetString("SLYTHE_HOST"), &http.Client{})
	_, err := client.UpdateCPUAutoscalingPolicy(context.Background(), cp)
	return err
}
