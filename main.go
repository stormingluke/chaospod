package main

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	emptyString string = ""
)

// initBuildLogger constructs a zap logger with the ECS format. This format is easily picked up and parsed by Elasticsearch compatible APIs and other ingestors/parsers.
// it is returned here a Logger but the functions that use this logger all use Sugared implementations.
func initBuildLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.Level.SetLevel(zapcore.InfoLevel)
	// setting the elasticsearch encoding here.
	encoderConfig := ecszap.NewDefaultEncoderConfig()
	// setting info level for logging clarity.
	core := ecszap.NewCore(encoderConfig, os.Stdout, zap.InfoLevel)
	return zap.New(core, zap.AddCaller())
}

// initBuildKubeClient creates an incluster config client. If it fails to do so it logs an error and calls os.Exit(1).
// There is no need to continue the process if there is no kube clientset to use.
func initBuildKubeClient(log *zap.Logger) (client *kubernetes.Clientset) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Sugar().Fatalf("failed to create inclusterConfig with error: %v", err)
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Sugar().Fatalf("failed to create kube clientset with error: %v", err)
	}
	return clientset
}

// main is the entrypoint it calls the two init functions to set up a logger using zap and a kubeclient set.
// These are then passed as values to subsequent functions and are not used in a struct as methods: based on Bill Kennedy's notebook from Arden Labs.
func main() {
	logger := initBuildLogger()
	kubeClient := initBuildKubeClient(logger)
	namespace, timeout := utilGetEnvVars(logger)
	// pullTheLeverKronk takes a logger and a kubeclient and starts the process of deleting a pod.
	// if it fails then the application is shut down.
	err := pullTheLeverKronk(logger, kubeClient, namespace, timeout)
	if err != nil {
		log.Default().Fatalf("wrong lever! with error: %v exiting", err)
	}
	// scaleDownChaos runs if the previous function ran successfully.
	// NOTE I call fatal in the previous function as a pod selected for deletion may still be starting up.
	// Kubernetes' state engine will restart this deployment when it is exited because it has not been
	// scaled down to 0 yet.
	err = scaleDownChaos(logger, kubeClient, namespace)
	if err != nil {
		log.Default().Fatalf("wrong lever! with error: %v exiting", err)
	}
}

// pullTheLeverKronk deletes a random pod from the targeted namespace. It is also a reference to The Emperor's New Groove.
func pullTheLeverKronk(log *zap.Logger, kubeClient *kubernetes.Clientset, namespace string, timeout int) (err error) {
	cctx, deadline := context.WithDeadline(context.Background(), time.Now().Add(time.Duration(timeout)*time.Second))
	defer deadline()
	// only expect 1 results to be passed over the channel, using buffer of 1, and adheres to Ubers excellent go style guide.
	// the channel is buffered, so the 'send' in the goroutine is nonblocking. This is a common pattern to prevent goroutine leaks in case the channel is never read from: https://gobyexample.com/timeouts
	// there is no need to return another type over the channel. The subsequent operations will output the results of the function's invocation.
	done := make(chan bool, 1)
	go deletePod(log, kubeClient, namespace, done)
	select {
	case ok := <-done:
		if !ok {
			return errors.New("delete action on pod failed")
		} else {
			// declare an error inside the select statement to prevent a nill pointer to an uninitialsed error. No error type is
			// returned from the function and the return declaration is also uninitialsed. Creating a *new() value allocates a nill
			// error in memory allowing the function to complete with a err == nil value.
			return nil
		}
	// context is cancelled after 10 seconds to prevent unnecessary cloud time
	case <-cctx.Done():
		return cctx.Err()
	}
}

// deletePod uses two functions to delete a pod. There is a good example including a retry peration on a pod in the release-1.18 example at:
// https://github.com/kubernetes/client-go/tree/master/examples/create-update-delete-deployment
func deletePod(log *zap.Logger, k *kubernetes.Clientset, namespace string, done chan bool) {
	ctx := context.Background()

	pods, err := utilListPodsInNamespace(k, ctx, namespace)
	if err != nil {
		log.Sugar().Errorf("failed find pods in the namespace: %s with error: %v", namespace, err)
		done <- false
	}
	podToDelete, err := utilRandomPod(pods)
	if err != nil {
		log.Sugar().Errorf("failed to select a pod to delete: %v", err)
		done <- false
	}
	// DeleteOptions can be a relevant selection criteria. Specifically the PropogationPolicy which can chain the delete action.
	err = k.CoreV1().Pods(namespace).Delete(ctx, podToDelete, metav1.DeleteOptions{})
	if err != nil {
		log.Sugar().Errorf("failed to delete pod: %s with error: %v", podToDelete, err)
		done <- false
	}
	done <- true
}

// scaleDownChaos is invoked by main after a random pod has been deleted or an attempt has been made to do so.
// It scales its own deployment to 0 thereby preventing any crashloops or excessive deletion actions.
func scaleDownChaos(log *zap.Logger, k *kubernetes.Clientset, namespace string) (err error) {
	ctx := context.Background()
	currentScale, err := k.AppsV1().Deployments(namespace).GetScale(ctx, "chaospod", metav1.GetOptions{})
	if err != nil {
		log.Sugar().Errorf("current scale of deployment unkown with error: %v", err)
	}

	// here we set the downscale variable to the pointer in memory to what the current deployment scale is and
	// replace it with 0. This is a one way function, bringing the deployment up again will also trigger this
	// downscaling and will bring any number of deployments down to 0 once the first run has completed.
	downscale := *currentScale
	downscale.Spec.Replicas = 0
	_, err = k.AppsV1().Deployments(namespace).UpdateScale(ctx, "chaospod", &downscale, metav1.UpdateOptions{})
	if err != nil {
		log.Sugar().Errorf("failed to scale the chaos deployment down with error: %v", err)
		return err
	}
	return nil
}
