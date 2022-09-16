package main

import (
	"context"
	"path/filepath"
	"testing"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type testKubeClient struct {
	logger     *zap.Logger
	restClient *kubernetes.Clientset
	namespace  string
	podList    *corev1.PodList
}

var testClient testKubeClient

func initTestKubeClient() {
	log := initBuildLogger()
	// taken from https://github.com/kubernetes/client-go/blob/master/examples/out-of-cluster-client-configuration/main.go
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = 
			filepath.Join(home, ".kube", "config")
	} 
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
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

func init() {
		initTestKubeClient()
}

func Test_pullTheLeverKronk(t *testing.T) {
	type args struct {
		log        *zap.Logger
		kubeClient *kubernetes.Clientset
		namespace  string
		podName    string
		timeout    int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
				name: "Test1: delete random pod",
			args: args{
				log:        testClient.logger,
				kubeClient: testClient.restClient,
				namespace:  testClient.namespace,
				// podName is emptyString by default but can be a specific target
				podName:    "",
				timeout:    10,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pullTheLeverKronk(tt.args.log, tt.args.kubeClient, tt.args.namespace, tt.args.podName, tt.args.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("pullTheLeverKronk() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_scaleDownChaos(t *testing.T) {
	type args struct {
		log       *zap.Logger
		k         *kubernetes.Clientset
		namespace string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
			{
					name:    "Test1: scale down deployment",
				args:    args{
					log:       testClient.logger,
					k:         testClient.restClient,
					namespace: testClient.namespace,
				},
				wantErr: false,
			},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := scaleDownChaos(tt.args.log, tt.args.k, tt.args.namespace); (err != nil) != tt.wantErr {
				t.Errorf("scaleDownChaos() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

