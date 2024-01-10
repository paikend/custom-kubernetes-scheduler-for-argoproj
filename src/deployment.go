package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	clientset, err2 = kubernetes.NewForConfig(config)
	api             = clientset.CoreV1()
)

func ProcessDeployment(nameSpace string, deploymentName string, serviceInstanceNum int, flow string) (map[string]string, bool) {
	result := true
	nodeselectors := make(map[string]string)
	deploymentAnnotations := map[string]string{}

	deploymentsClient := clientset.AppsV1().Deployments(nameSpace)
	deploymentData, getErr := deploymentsClient.Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if getErr != nil {
		panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
	}

	deploymentAnnotations = deploymentData.Annotations
	if strategy, ok := deploymentAnnotations["custom-pod-schedule-strategy"]; ok && strategy != "" {
		nodeselectors, result = ProcessData(strategy, deploymentData, deploymentName, nameSpace, serviceInstanceNum, flow)
	}
	return nodeselectors, result
}
