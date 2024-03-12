# Application Operator

The Application Operator is a Kubernetes operator built using the Go plugin of Operator SDK. It automates the process of deploying an application as a Kubernetes Deployment and exposing it as a Kubernetes Service.

## Prerequisites

- Kubernetes cluster
- `kubectl` installed and configured to interact with your cluster
- Operator SDK

## Operator Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/application-operator.git
cd application-operator
```

2. Install the CRD:

```bash
kubectl apply -f config/crd/bases/
```

3. Install the Operator:

```bash
kubectl apply -f config/manager/
```

## Usage

1. Create an Application resource:

```yaml
apiVersion: api.app.op/v1alpha1
kind: Application
metadata:
  labels:
    app.kubernetes.io/name: application
    app.kubernetes.io/instance: application-sample
    app.kubernetes.io/part-of: application-operator
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/created-by: application-operator
  name: application-sample
spec:
  replicas: 2
  image: sayedimran/fastapi-sample-app:v4
  port: 7000
  env:
  - name: APP_ENV
    value: production
  - name: APP_NAME
    value: fastapi-sample-app
```

2. Apply the resource:

```bash
kubectl apply -f config/samples/
```

3. Verify that the Application resource has been created:

```bash
kubectl get app
```

4. Verify that the Deployment and Service resources have been created:

```bash
kubectl get deployment
kubectl get service
```
