# Tyk operator

Tyk Operator brings Full Lifecycle API Management capabilities to Kubernetes.
Configure Ingress, APIs, Security Policies, Authentication, Authorization, Mediation and more - all using GitOps best practices with Custom Resources and Kubernetes-native primitives.

## Usage

```
helm repo add tyk-operator https://tyktechnologies.github.io/tyk-operator
helm repo update
```

Then install
```
helm install ci tyk-operator/tyk-operator
```