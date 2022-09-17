/*
This util.go file contains some utility functions that are used for:
- getting variables out of the environment if any are set.
- randomly selecting a pod from a list
-
The test files were generated using a go utility binary.
*/
package main

import (
	"context"
	"os"
	"strings"
	"testing"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

func Test_utilRandomPod(t *testing.T) {
	type args struct {
		pods *corev1.PodList
	}
	tests := []struct {
		name            string
		args            args
		wantPodToDelete string
		wantErr         bool
	}{
		{
			name: "Test01: return random pod",
			args: args{
				pods: testClient.podList,
			},
			wantPodToDelete: "nginx",
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPodToDelete, err := utilRandomPod(tt.args.pods)
			if (err != nil) != tt.wantErr {
				t.Errorf("utilRandomPod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(gotPodToDelete, tt.wantPodToDelete) {
				t.Errorf("utilRandomPod() got = %v, want %v", gotPodToDelete, tt.wantPodToDelete)
			}
		})
	}
}

func Test_utilListPodsInNamespace(t *testing.T) {
	type args struct {
		k         *kubernetes.Clientset
		ctx       context.Context
		namespace string
	}
	tests := []struct {
		name              string
		args              args
		// this is a flaky test as the number of pods is dependent on the speed in which the scaleDownDeployment test runs. 
		// if the downscale of the deployments takes longer than it does to get to- and run this test then there are 4 pods in
		// the test environment. Otherwise there are 3.
		wantPodListLength map[int]string
		wantErr           bool
	}{
		{
			name: "Test01: check all returned pods",
			args: args{
				k:         testClient.restClient,
				ctx:       context.Background(),
				namespace: testClient.namespace,
			},
			wantPodListLength: map[int]string{3:"noChaosPod", 4:"withChaosPod"},
			wantErr:           false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPodList, err := utilListPodsInNamespace(tt.args.k, tt.args.ctx, tt.args.namespace)
			if (err != nil) != tt.wantErr {
				t.Errorf("utilListPodsInNamespace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var numberOfPods []string
			for _, pod := range gotPodList.Items {
				numberOfPods = append(numberOfPods, pod.Name)	
			}
			if _, ok := tt.wantPodListLength[len(numberOfPods)]; !ok {
					t.Errorf("got = %v want pod length 3 or 4", len(numberOfPods))
			}

		})
	}
}

func Test_utilGetEnvVars(t *testing.T) {
	os.Setenv("TARGET_NAMESPACE", "workloads-test")
	os.Setenv("TIMEOUT", "42")

	type args struct {
		log *zap.Logger
	}
	tests := []struct {
		name          string
		args          args
		wantNamespace string
		wantPodName   string
		wantTimeout   int
	}{
		{
			name: "Test01: check set vars and no defaults",
			args: args{
				log: initBuildLogger(),
			},
			wantNamespace: "workloads-test",
			wantPodName:   "nginx-test",
			wantTimeout:   42,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNamespace, gotTimeout := utilGetEnvVars(tt.args.log)
			if gotNamespace != tt.wantNamespace {
				t.Errorf("utilGetEnvVars() gotNamespace = %v, want %v", gotNamespace, tt.wantNamespace)
			}
			if gotTimeout != tt.wantTimeout {
				t.Errorf("utilGetEnvVars() gotTimeout = %v, want %v", gotTimeout, tt.wantTimeout)
			}
		})
	}
	os.Unsetenv("TARGET_NAMESPACE")
	os.Unsetenv("TIMEOUT")
}
