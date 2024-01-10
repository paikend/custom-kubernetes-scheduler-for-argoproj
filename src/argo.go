package main

import (
	"context"
	"fmt"
	argoclientset "github.com/argoproj/argo-rollouts/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	argoClientset, err3 = argoclientset.NewForConfig(config)
)

func ProcessRollout(nameSpace string, rolloutName string, serviceInstanceNum int, flow string) (map[string]string, bool) {
	result := true
	nodeselectors := make(map[string]string)
	rolloutAnnotations := map[string]string{}
	rolloutClient := argoClientset.ArgoprojV1alpha1().Rollouts(nameSpace)
	rolloutData, getErr := rolloutClient.Get(context.TODO(), rolloutName, metav1.GetOptions{})
	if getErr != nil {
		panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
	}
	rolloutAnnotations = rolloutData.Annotations
	if strategy, ok := rolloutAnnotations["custom-pod-schedule-strategy"]; ok && strategy != "" {
		nodeselectors, result = ProcessData(strategy, rolloutData, rolloutName, nameSpace, serviceInstanceNum, flow)
	}
	return nodeselectors, result
}
