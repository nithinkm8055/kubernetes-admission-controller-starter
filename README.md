# Kubernetes Admission Controller Webhook
This repository implements a mutating webhook that removes provided labels from Namespaces.

The webhook runs as an HTTPS server and is triggered on `CREATE` action of a `Namespace` Object.

## Build the application

```sh
make go-build
```

## Build and Tag Application Docker Image

```sh
make docker-image
```

# Deploy MutatingWebhookConfiguration

Create a kubernetes secret containing the tls certificate and key for webhook-server. This is mounted to the kubernetes deployment.

The `deployment/` directory contains the kubernetes manifest to deploy the webhook and mutating webhook configuration.

The `caBundle` field in MutatingWebhookConfiguration needs to be replaced with a base64 encoded CA certificate for the webhook server.

# Tests

The admission controller can be tested using files in `samples/` directory.
1. `kustomization.yaml` file provides an example of how the application is patched using a JSON6902 Patch. Similar JSONPatch is applied by the webhook to mutate labels.

2. `sampleRequestBody.json` is a sample request format containing `AdmissionReview` object received by the Webhook Server.

```
➜  kubectl logs -f -n webhook-demo webhook-server-54d688b94b-nnnhm 
2023/12/28 14:38:07 Handling webhook request, traceID: a8e2f2a5-f47e-42a4-8737-3113f6ec93e1
2023/12/28 14:38:07 Webhook request handled successfully, traceID: a8e2f2a5-f47e-42a4-8737-3113f6ec93e1
```

## References
> [A Guide to Kubernetes Admission Controllers](https://kubernetes.io/blog/2019/03/21/a-guide-to-kubernetes-admission-controllers/)