# cloud-native-demo
A demo for the power of cloud native application design


## Installation

Requirements:

* a Kubernetes Cluster to which you have writing API access. The cluster should have the following services installed:

  * Linkerd 2.x (minimal version: `stable-2.7.x` or `edge 2024.11.8`)
  * grafana-operator
  * prometheus

  If you do not have access to a managed K8s offering, that provides these components, try the [K8s Test Rig](https://github.com/cypherfox/k8s-test-rig), which also provides

  * cert-manager
  * vault
  * external-secrets
  * external-dns
  * ingress-nginx

* the [step-cli](https://github.com/smallstep/cli) locally installed.
* the [helm v3.x](https://helm.sh/docs/intro/install/) binary locally installed.

TODO: make this into a shell script and execute it via a docker image.

Then execute the following script:

``` bash
step certificate create root.linkerd.cluster.local ca.crt ca.key \
--profile root-ca --no-password --insecure

step certificate create identity.linkerd.cluster.local issuer.crt issuer.key \
--profile intermediate-ca --not-after 8760h --no-password --insecure \
--ca ca.crt --ca-key ca.key

pushd deploy/helm/cloud-native-demo && helm dependency update && popd

# Linkerd installed as a pre-requisite. See note below.
# Skip this step if your kubernetes cluster comes with Linkerd v2.x pre-installed.

helm repo add linkerd-edge https://helm.linkerd.io/edge && helm repo update linkerd

helm upgrade linkerd-crds linkerd-edge/linkerd-crds \
  --install --version 2024.11.8\
  -n linkerd --create-namespace

helm upgrade linkerd-control-plane \
  --install --version 2024.11.8 \
  -n linkerd \
  --set-file identityTrustAnchorsPEM=ca.crt \
  --set-file identity.issuer.tls.crtPEM=issuer.crt \
  --set-file identity.issuer.tls.keyPEM=issuer.key \
  linkerd-edge/linkerd-control-plane

#
# 
#
helm upgrade cloud-native-demo \
  --install \
  -n cloud-native-demo --create-namespace \
  ./deploy/helm/cloud-native-demo
```

## Linkerd Installation

Due to the business decisions by Buoyant, stable releases of Linkerd 2.x are only
available by purchasing a license. Therefore this is not a chart dependency of the
cloud native demo anymore. The script above uses a relatively current edge release.


## Quick Check

``` bash
kubectl proxy

```
Web-Browser to: 

* Bug Simulator: http://127.0.0.1:8001/api/v1/namespaces/cloud-native-demo/services/cloud-native-demo:80/proxy/
* Emojivoto Web Frontend: http://127.0.0.1:8001/api/v1/namespaces/cloud-native-demo/services/cloud-native-demo-emojivoto-web-svc:80/proxy/