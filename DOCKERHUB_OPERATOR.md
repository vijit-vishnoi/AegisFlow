# AegisFlow Operator

Kubernetes operator for AegisFlow custom resources.

Operator reconciles gateway, provider, route, tenant, and policy resources. Use matching operator and gateway versions.

## Image

```text
saivedant169/aegisflow-operator:0.9.0
```

Same image publishes at `ghcr.io/saivedant169/aegisflow-operator:0.9.0`.

## Install manifests

Deployment and RBAC manifests live under [`deployments/operator`](https://github.com/saivedant169/AegisFlow/tree/main/deployments/operator).

Review namespace, service account, RBAC, image tag, webhook certificates, and custom resource definitions before applying them.

Project documentation: [github.com/saivedant169/AegisFlow](https://github.com/saivedant169/AegisFlow)
