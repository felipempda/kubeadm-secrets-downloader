# kubeadm-secrets-downloader

Simple Go script to decrypt the certificates in secret created by kubeadm during ```upload-certs``` phase. This is for test purposes only, for folks interested in learning kubeadm internals.

# Usage

```bash
go build
./kubeadm-secrets-downloader --help
```

# (optional) Upload-certs phase

 If there is no secret in Kubernetes created by kubeadm or if it's already deleted you can genereate it again (this will not generate new certs, dont' worry).

Run this from a control plane node:

```bash
sudo kubeadm init phase upload-certs --upload-certs
```

As you can see a new secret is created in Kubernetes:

```yaml
[upload-certs] Storing the certificates in Secret "kubeadm-certs" in the "kube-system" Namespace
[upload-certs] Using certificate key:
a7de21ae2185df1b0dd1411ed079ee43dabc092ebfd852d17a4074afd838f725
```

You can view the name of the certs inside of this secret:
```bash
kubectl describe secrets kubeadm-certs -n kube-system
```

```yaml
Name:         kubeadm-certs
Namespace:    kube-system
Labels:       <none>
Annotations:  <none>

Type:  Opaque

Data
====
ca.crt:              1127 bytes
ca.key:              1707 bytes
etcd-ca.crt:         1114 bytes
etcd-ca.key:         1707 bytes
front-proxy-ca.crt:  1143 bytes
front-proxy-ca.key:  1707 bytes
sa.key:              1707 bytes
sa.pub:              479 bytes
```

If you happen to check the contents of this secret you will see each certificate is encrypted:

```bash
kubectl get secrets kubeadm-certs -n kube-system -o yaml
```

```yaml
apiVersion: v1
data:
  ca.crt: a2+8qXAmF3oMnwusbMXaV/TaN2ChbwCFJSlIDUcEFSib+5WUH2FhI8KkPVuRVoOGUbU9n9PgRh8plA2xkXNvKGOy555h+vuGh2NzI1Po++lqhN8ISwvlh0bKxkYsrLUxxwDFe6v2FfFLoDjPnIwD599IHr5qzSd2LiPb6V9Glg9sFhW
  (...)
  +omSaamQQ3tLChv0UVRStCv69soGhERaFy7WFxSK7bv5D1hmc1F6yrx6otSVrwFIIdA3u8h/dvo=
  ca.key: WkJLACBnVJ1x09olHpn5lnEHxps0cCQDqpJzC2Nml/+URcMXDoeNQqI+i1lsLBQk1W+tOydOp2KsElhGjxpPdPCbXE4vr0dI4wbgXFjBBDsZBwz2hT7s5IzREDS/
  (...)
DzRNUYt9nfa6fyxaJWWRtArTeKBbLnYcBIUvvJRiRqj61PzqLg8rLCRCV3IRfItC4YxE+Bp+cJDP5jmX4lePuCSHuEMGoCSfY5O8w4it432kEPPY9A8AF0dymFXHR6MeXetmOaVmfC3zG40tjE60+nw
  (...)
```

This is not only base64 encoded but also encrypted by the certificate key showed previously.

# Example

Once  you have the certificate key you can use it to see the contents of certificates with this script (using --key argument). 

```yaml
./kubeadm-secrets-downloader --key "a7de21ae2185df1b0dd1411ed079ee43dabc092ebfd852d17a4074afd838f725"
sa.pub :
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArdc/GuXLA2OHm+quFaAf
P8hulvAd0P9u83xyX9CsZwPGQbRMQFe1wk1RwlUU60QNCmZSJSdXV9my5Z77fH/L
oODT0oMj...
7wIDAQAB
-----END PUBLIC KEY-----


ca.crt :
-----BEGIN CERTIFICATE-----
MIIC/jCCAeagAwIBAgIBADANBgkqhkiG9w0BAQsFADAVMRMwEQYDVQQDEwprdWJl
cm5
1ovxFK0v7Wg9TyJS2SKtinBjPV8cWVpYkFFbQ6UbgrZpLqu2b3B39AoCEaRpL7WR
visqXJ8fKIgrzwN1J2vY4RK+ptTXpag4w6iLqKPuxRO5aYUfuhzr2rzW+GKSQ3Nv
1rY=
-----END CERTIFICATE-----

(...)
```

Of you course there is no need for that if you have access to the control plane directly. This is just to understand how kudeadm uses Kubernetes Secrets as a way to securely storage control plane certificates in order to bootstrap new nodes in a easy way.