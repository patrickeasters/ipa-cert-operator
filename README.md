# FreeIPA Cert Operator

## Overview

This operator automates the provisioning and renewal of certificates from
FreeIPA. This allows OpenShift users to request and manage valid certificates
just as they would any other resources.

## Usage

The operator handles issuing both server and client certificates. Generated
certificates are stored in a secret, except in the case of server certificates
requested for routes.

### Assumptions

* A FreeIPA (or Red Hat Identity Management) server has already been configured
* A service account with permission to read/request certificates is configured
  in FreeIPA
* Host or user records already exist in FreeIPA for any certificate subjects

### Host Certificates
Host certificates can be generated for any host principal in FreeIPA. Host
records must already exist for both the `cn` and all SANs. A secret will be
generated with the pattern `$name-tls` (e.g. `myservice-tls`).

```yaml
apiVersion: cert.patrickeasters.com/v1alpha1
kind: IpaCert
metadata:
  name: myservice
spec:
  type: host
  cn: myservice.lab.easte.rs
  AdditionalNames:
    - myservice.easte.rs
```

### Host Certificates (routes)
You can also use this operator to generate certificates for edge or re-encrypt
routes. This is accomplished by simply adding an annotation to the route. No
secret will be created since the certificate and private key are stored in the
`TLSConfig` of the route spec.

```yaml
apiVersion: v1
kind: Route
metadata:
  name: myroute
  annotations:
    cert.patrickeasters.com/ipa-managed: 'true'
spec:
  host: myroute.easte.rs
  tls:
    termination: reencrypt
```

### User Certificates

Client certificates can be generated for a user principal in FreeIPA. Like host
certificates, a secret will be generated with the pattern `$name-tls`.

```yaml
apiVersion: cert.patrickeasters.com/v1alpha1
kind: IpaCert
metadata:
  name: sa1
spec:
  type: user
  cn: sa1
```

## Deploying

### Prerequisites

- [oc][oc_tool] version v3.9.0+.
- Access to an OpenShift v3.9.0+ cluster.

### Steps

Configure RBAC for the ipa-cert-operator and its related resources:

```sh
oc create -f deploy/service_account.yaml
oc create -f deploy/role.yaml
oc create -f deploy/role_binding.yaml
```

Configure the IpaCert CRD:

```sh
oc create -f deploy/cert_v1alpha1_crd.yaml
```

Update config and credentials:

```sh
cp deploy/config/config_example.yaml deploy/config/config.yaml
cp deploy/config/credentials_example.yaml deploy/config/credentials.yaml
# Edit the values in your editor of choice
atom deploy/config/config.yaml
atom deploy/config/credentials.yaml
```

Deploy ipa-cert-operator:

```sh
oc create -f deploy/config/config.yaml
oc create -f deploy/config/credentials.yaml
oc create -f deploy/operator.yaml
```

## Developing/Building

### Prerequisites

- [dep][dep_tool] version v0.5.0+.
- [go][go_tool] version v1.10+.
- [docker][docker_tool] version 17.03+.
- [oc][oc_tool] version v3.9.0+.
- Access to an OpenShift v3.9.0+ cluster.

### Install the Operator SDK CLI

First, checkout and install the operator-sdk CLI:

```sh
mkdir -p $GOPATH/src/github.com/operator-framework
cd $GOPATH/src/github.com/operator-framework
git clone https://github.com/operator-framework/operator-sdk
cd operator-sdk
git checkout master
make dep
make install
```

### Build Steps

Checkout this repository:

```sh
mkdir $GOPATH/src/github.com/patrickeasters
cd $GOPATH/src/github.com/patrickeasters
git clone https://github.com/patrickeasters/ipa-cert-operator.git
cd ipa-cert-operator
```

Vendor the dependencies:

```sh
dep ensure
```

### Run the operator locally
Since the operator does require some config, we'll export some environment vars
and run it using the `operator-sdk` CLI. There are sane defaults for most
options, but at minimum, the IPA connection settings must be configured.

```sh
export IPA_HOST=ipa.lab.easte.rs
export IPA_REALM=LAB.EASTE.RS
export IPA_USER=admin
export IPA_PASSWORD=password
export OPERATOR_NAME=ipa-cert-operator
operator-sdk up local --namespace=default
```

### Build the operator

Build the operator image and push it to a public registry such as quay.io:

```sh
export IMAGE=quay.io/patrickeasters/ipa-cert-operator:v0.0.1
operator-sdk build $IMAGE
docker push $IMAGE
```


[client_go]:https://github.com/kubernetes/client-go
[operator_sdk]:https://github.com/operator-framework/operator-sdk
[dep_tool]:https://golang.github.io/dep/docs/installation.html
[go_tool]:https://golang.org/dl/
[docker_tool]:https://docs.docker.com/install/
[oc_tool]:https://github.com/openshift/origin/releases/
