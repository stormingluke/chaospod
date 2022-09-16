/*
This util.go file contains some utility functions that are used for:
- getting variables out of the environment if any are set.
- randomly selecting a pod from a list
-
The test files were generated using a go utility binary.
*/package main

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type testKubeClient struct {
	restClient *kubernetes.Clientset
	namespace  string
	podList    *corev1.PodList
}

var testClient *testKubeClient

func initTestKubeClient() {
	log := initBuildLogger()
	// copied from https://github.com/kubernetes/client-go/blob/master/examples/out-of-cluster-client-configuration/main.go
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String(
			"kubeconfig",
			filepath.Join(home, ".kube", "config"),
			"(optional) absolute path to the kubeconfig file",
		)
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Sugar().Fatalf("failed to create kubeconfig from flags with error: %v", err)
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Sugar().Fatalf("failed to create kubeconfig for tests with error: %v", err)
	}
	testClient.restClient = clientset
	testClient.namespace = "workloads"
	podsList, err := utilListPodsInNamespace(clientset, context.Background(), testClient.namespace)
	if err != nil {
		log.Sugar().Fatalf("failed to get pod list with error: %v", err)
	}
	testClient.podList = podsList
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
			if strings.Contains(gotPodToDelete, tt.wantPodToDelete) {
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
			{
					name:        "Test01: check all returned pods",
				args:        args{
					k:         testClient.restClient,
					ctx:       context.Background(),
					namespace: testClient.namespace,
				},
				wantPodList: testClient.podList,
				wantErr:     false,
			},
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
	os.Setenv("TARGET_NAMESPACE", "workloads-test")
	os.Setenv("WRONG_LEVER", "nginx-test")
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
					name:          "Test01: check set vars and no defaults",
				args:          args{
					log: initBuildLogger(),
				},
				wantNamespace: "workloads-test",
				wantPodName:   "nginx-test",
				wantTimeout:   42,
			},
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
	os.Unsetenv("TARGET_NAMESPACE")
	os.Unsetenv("WRONG_LEVER")
	os.Unsetenv("TIMEOUT")
}

