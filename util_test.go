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
	"reflect"
	"testing"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type testKubeClient struct {
		restClient *kubernetes.Clientset
		namespace string
		podList *corev1.PodList
}

func initBuildTestKubeClient(log *zap.SugaredLogger) {

}

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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPodToDelete, err := utilRandomPod(tt.args.pods)
			if (err != nil) != tt.wantErr {
				t.Errorf("utilRandomPod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPodToDelete != tt.wantPodToDelete {
				t.Errorf("utilRandomPod() = %v, want %v", gotPodToDelete, tt.wantPodToDelete)
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
		name        string
		args        args
		wantPodList *corev1.PodList
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPodList, err := utilListPodsInNamespace(tt.args.k, tt.args.ctx, tt.args.namespace)
			if (err != nil) != tt.wantErr {
				t.Errorf("utilListPodsInNamespace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPodList, tt.wantPodList) {
				t.Errorf("utilListPodsInNamespace() = %v, want %v", gotPodList, tt.wantPodList)
			}
		})
	}
}

func Test_utilGetEnvVars(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNamespace, gotPodName, gotTimeout := utilGetEnvVars(tt.args.log)
			if gotNamespace != tt.wantNamespace {
				t.Errorf("utilGetEnvVars() gotNamespace = %v, want %v", gotNamespace, tt.wantNamespace)
			}
			if gotPodName != tt.wantPodName {
				t.Errorf("utilGetEnvVars() gotPodName = %v, want %v", gotPodName, tt.wantPodName)
			}
			if gotTimeout != tt.wantTimeout {
				t.Errorf("utilGetEnvVars() gotTimeout = %v, want %v", gotTimeout, tt.wantTimeout)
			}
		})
	}
}

func Test_utilLookupEnvVar(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name         string
		args         args
		wantEnvValue string
		wantErr      bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEnvValue, err := utilLookupEnvVar(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("utilLookupEnvVar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotEnvValue != tt.wantEnvValue {
				t.Errorf("utilLookupEnvVar() = %v, want %v", gotEnvValue, tt.wantEnvValue)
			}
		})
	}
}
