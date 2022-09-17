/*
This util.go file contains some utility functions that are used for:
- getting variables out of the environment if any are set.
- randomly selecting a pod from a list
-
*/
package main

import (
	"context"
	"errors"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// utilRandomPod first filters out the chaospod that has been deployed. It creates an array of podNames which are then passed into
// a second loop that randomly selects a pod based on the index of the array. If nothing is found then an emptyString is returned
// and an error.
func utilRandomPod(pods *corev1.PodList) (podToDelete string, err error) {
	var filteredPodList []string
	rand.Seed(time.Now().UnixNano())

	for _, filteredPod := range pods.Items {
		if !strings.Contains(filteredPod.Name, "chaospod") {
			filteredPodList = append(filteredPodList, filteredPod.Name)
		}
	}
	return filteredPodList[rand.Intn(len(filteredPodList))], nil
}

// listPodsInNamespace lists all the pods in the namespace and returns them in a list of v1.PodList.
// PS I do not understand why gofmt is turning this funciton into a multiline declaration.
func utilListPodsInNamespace(
	k *kubernetes.Clientset,
	ctx context.Context,
	namespace string,
) (podList *corev1.PodList, err error) {
	// get pods in all the namespaces by omitting namespace
	pods, err := k.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return &corev1.PodList{}, errors.New("could not get pods")
	}
	return pods, nil
}


// utilGetEnvVars gets a target namespace from the environment varibale TARGET_NAMESPACE if it is not set it defaults to 'workloads'.
// it can accept a podName through the WRONG_LEVER environment variable but defaults to an emptyString if not set.
// It looks for a TIMEOUT environment variable and defaults to a 10 Second operation if no timeout is specified.
func utilGetEnvVars(log *zap.Logger) (namespace string, timeout int) {
	namespace, err := utilLookupEnvVar("TARGET_NAMESPACE")
	if err != nil {
		log.Sugar().Infof("did not receive a TARGET_NAMESPACE, defaulting to the 'workloads' namespace: %v", err)
		namespace = "workloads"
	}
	timeoutVar, err := utilLookupEnvVar("TIMEOUT")
	if err != nil {
		log.Sugar().Infof("did not receive a timout variable, defaulting to 10 Seconds: %v", err)
		timeout = 10
	} else {
		timeout, err = strconv.Atoi(timeoutVar)
		if err != nil {
			log.Sugar().Error("failed to convert to valid int, defaulting to 10")
			timeout = 10
		}
	}
	return namespace, timeout
}

// utilLookupEnvVar is a wrapper around the os.LookupEnv function and returns an error if the environment variable is not set.
func utilLookupEnvVar(key string) (envValue string, err error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return emptyString, errors.New("no key found for supplied value")
	}
	return val, nil
}
