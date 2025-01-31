# Operator context

This is a resource used to define values for communicating with admin api of
tyk ce/pro deployment.

`OperatorContext` can be referenced on all custom resources by `contextRef` property.
When a custom resource is applied with `contextRef` set then all the operations
conducted by the operator will use the `OperatorContext` supplied by the `contextRef` 
to perform reconciliation.

# Defining OperatorContext

Annotated `OperatorContext` for a community edition deployment looks like this

```yaml
apiVersion: tyk.tyk.io/v1alpha1
kind: OperatorContext
metadata:
  name: community-edition
spec:
  env:
    # The mode of the admin api
    # ce - community edition (open source gateway)
    # pro - dashboard (requires a license)
    mode: ce
    # The authorization token this will be set in x-tyk-authorization header on the
    # client while talking to the admin api
    auth: foo
    # Org ID to use
    org: myorg
    # The url to the tyk open source gateway deployment admin api
    url: http://tyk.tykce-control-plane.svc.cluster.local:8001
    # Set this to true if you want to skip tls certificate and host name verification
    # this should only be used in testing
    insecureSkipVerify: true
    # For ingress the operator creates and manages ApiDefinition resources, use this to configure
    # which ports the ApiDefinition resources managed by the ingress controller binds to.
    # Use this to override default ingress http and https port
    ingress:
      httpPort: 8000
      httpsPort: 8443
```

# Using secret for sensitive information

Whilst it is possible to set `.spec.env.auth` directly in the `OperatorContext` object, better security can be achieved by replacing sensitive data with values contained within a referenced secret.
Create a `k8s` secret `tyk-operator-system` with our sensitive info for auth and values

```sh
kubectl create secret -n tyk-operator-system generic tyk-operator-conf \
  --from-literal "TYK_AUTH=foo" \
  --from-literal "TYK_ORG=myorg" \
```

We can now reference our secret in the `OperatorContext` resource with `.spec.secretRef`

```yaml
apiVersion: tyk.tyk.io/v1alpha1
kind: OperatorContext
metadata:
  name: community-edition
spec:
  secretRef:
    namespace: tyk-operator-system
    name: tyk-operator-conf
  env:
    mode: ce
    url: http://tyk.tykce-control-plane.svc.cluster.local:8001
    insecureSkipVerify: true
    ingress:
      httpPort: 8000
      httpsPort: 8443
```

Mappings between `.spec.env` properties and secret `.spec.data` keys

| secret key | .spec.env |
|------------|-----------|
| TYK_MODE   | mode      |
| TYK_URL    | url       |
| TYK_AUTH   | auth      |
| TYK_ORG | org |
| TYK_TLS_INSECURE_SKIP_VERIFY | insecureSkipVerify |


# Referencing OperatorContext in ApiDefinion

We can refer  to the `OperatorContext` we created above to `ApiDefinion` resource using `context.ref` property like

```yaml
apiVersion: tyk.tyk.io/v1alpha1
kind: ApiDefinition
metadata:
  name: httpbin-with-context
spec:
  contextRef:
    name: community-edition
    namespace: default
  name: httpbin-with-context
  use_keyless: true
  protocol: http
  active: true
  proxy:
    target_url: http://httpbin.org
    listen_path: /httpbin-with-context
    strip_listen_path: true
```

Then `httpbin-with-context` api will be created on the gateway defined in  the `community-edition` `OperatorContext` resource.


# Default Context
The operator will not start without a valid context.

With  our helm charts we have
```
envFrom:
  - secretRef:
      name: tyk-operator-conf
```
That loads the default env from secret `tyk-operator-conf` . This might not be desired if you are already using operator context for resources.

One work around is to set dummy values that way the operator will start, and provide valid context when applying resources that will be used by the operator for reconciliation.

You can update your `values.yaml` to include

```yaml
envFrom:

envVars:
  - name: TYK_AUTH
    value: "xxx"
  - name: TYK_ORG
    value: "xxx"
  - name: TYK_MODE
    value: "pro"
  - name: TYK_URL
    value: "xxx"
```

This workaround will not be required in `v0.8.0` release