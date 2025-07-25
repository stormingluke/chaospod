# chaospod

Assignment for Control Plane

- [Summary](#summary)
  - [Why?](#why)
  - [Options](#options)
- [Requirements](#requirements)
  - [Quickstart](#quickstart)
  - [Run](#running-the-application)
- [Example](#example)
  - [Delete Random Pod inside Namespace](#delete-random-pod)
- [Considerations](#considerations)
  - [Quality](#quality)
  - [Permissions](#permissions)
  - [Testing](#testing)
  - [Build](#build)
  - [Monitoring](#monitoring)
- [Potential Next Steps](#next-steps)

# Summary

Kubernetes clusters drive towards a declared state. Since Kubernetes environments are
typically very dynamic, Pods are intended to be ephemeral, and could be destroyed at any
point.

## Why?

Since we should develop defensively, some companies test the resilience of workloads in advance by
doing performing delete/restart operations on purpose.

## Options

There are some configuration options available through Environment variables:

- TIMEOUT; defaults to 10 Seconds.
- TARGET_NAMESPACE; defaults to 'workloads'. Note that this is only relevant if it is tested against a namespace other than `workloads`. The Role/Rolebinding are configured to only work with the `workloads` namespace. Editting the namespace in the kustomization.yaml file will put these RBAC Objects in a different namespace.

# Requirements

This project currently demands the following requirements:

- GCP Account with Billing Enabled
- (Optional) Terraform.
- Artifact Store Enabled; there is a minimal terraform file and no GCP APIs are enabled.
- Cloudbuild; Capable of running builds and configuring triggers.
- A Kubernetes Cluster available; my cluster is a GKE Autopilot distribution.
- Kubectl with Kustomize installed
- Skaffold installed

## Quickstart

```sh
kubectl apply -f kustomize/testnginx/deployment.yaml
skaffold dev
```

**This action is (currently) configured to my own dev environment**
Skaffold dev will build and deploy the `podchaos` container.

## Run

To run this in on an actual cluster there's a few steps to go through because it's been configured to run in my GCP Project.

There are no dependencies other than a specific container image name that needs to be changed in the `kustomize/base/deployment.yaml` file and in the `skaffold.yaml` file.

The following chained command can be used to build the dockerfile using cloudbuild, create the kubectl manifests and deploy them.

```sh
skaffold build && skaffold render -p singlens --output render.yaml && skaffold apply render.yaml
```

# Example

There is an example situation/deployment specification in `kustomize/testnginx` this is a single file that contains 3 nginx pods with different names and a namespace.

## Delete Random Pod

After deploying the example project in `kustomize/testnginx` the command indicated above after editting the imageNames in the files will run and delete a random pod in the namespace. It will then scale down the chaos deployment to 0.

# Considerations

I've made a short list of things I've skipped or haven't configured properly.

## Quality

The repository contains a Github Acion pipeline in `.github/workflows/ci.yaml` that does 3 checks on the quality of the code:

- golangci; this is a popular linting tool that has been customised with the `.golangci.yaml` file.
- go vet; this is a built-in tool and detects any misallacted variables and dataraces among other common errors.
- go test; this is also a built-in tool, no external packages are used here. However, no kubernetes credentials are supplied to the repository pipeline.

Additionally the following two Github Applications are run:

- Snyk vulnerability scan; this tool checks to see if there are any known vulnerabilities in the currect codebase by checking dependency metadata against an online database of CVEs.
- dependabot; this tool is configured within github and is customised through the `.github/dependabot.yaml` file.

The last item in the list is a dependency specific to my repository.

## Permissions

I have given the chaos deployment excessive permissions with no particular sorting or grouping in the Role. The service account is attached to the deployment and its permissions are excessive.

## Testing

Given the 'soft' time limit I have not implemented any tests in the pipeline.

## Build

I'm building the app fresh each time in the Dockerfile with no real consideration for caching or build times. The `cloudbuild.yaml` uses the kaniko executor which is triggered on every push to this repository. However, the `skaffold.yaml` file does not use the `cloudbuild.yaml` and cannot access the same cache.
I also have not accounted for binary tags or any maintainership on the binary.
It is currently only a `cloudbuild.yaml` job, mostly because it's really short, easy and fast to configure and there's good interop with skaffold.

## Monitoring

Big One, this podchaos container checks to see if the target pod was deleted and then leaves it alone but it doesn't report this anywhere. Logging is minimal and I've not implemented any further metrics for observability.

# Potential (easy) Next Steps

The first step could be improving this codebase:

1. Integration testing: running Kind and Nginx in a pipeline would be a good start.
2. Tracking releases of the code to handle
3. Documenting the code properly.
4. Extract the hard-dependency on a kubernetes client version into a separate library that this codebase imports using the same API.
5. Verify the best practices around the Dockerfile, Kubernetes Manifests and Terraform Implementation.

Given that this container can be started and taken down quickly through skaffold there are some opportunities for potential next steps.

1. Perform the skaffold run step as a post-deployment pipeline step with a timeout
2. Expand the skaffold setup with profiles to determine configuration of the container: potentially creating cli-chaos on-demand

Focussing on the application itself:

1. The go chaos code could be extracted into a separate pkg and imported by other clients capable of interacting with the kube-api.
2. The application itself can be improved with logic to monitor state of a namespace during an outage. Implementing a 'Watcher' here on the namespace of the targeted pod could show how distributed transactions may fail.
3. By letting the application continue to monitor the Deployment or StatefulSet or other encapsulating object it could report on the timing between down-and-up to provide data towards reliability indicators.
4. Adding the container as a sidecar.
