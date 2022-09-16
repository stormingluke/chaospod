/* 
This container is a variation on the default go-client example code provided here: 
Rather than performing a GET action this container randomly selects a pod in the namespace and deletes it.
https://github.com/kubernetes/client-go/blob/master/examples/in-cluster-client-configuration/main.go

It is a single package application that consists of 2 files: one containing the core logic in the main.go file and the other
is the util.go file that contains some utility functions.

This application can be configured using 3 environment variables:
- TIMEOUT
- CONTAINER_NAME
- TARGET_NAMESPACE

If these are not set a util function will set some default values.

The application is started after deploying (at least) 1 replica. It will then randomly select a pod from inside the workloads namespace and delete it.
Then it will scale itself down to 0. 

If it fails to delete a pod or if it takes longer than the TIMEOUT variable (default 10 Seconds) then the application is also scaled down to 0.
*/
package main
