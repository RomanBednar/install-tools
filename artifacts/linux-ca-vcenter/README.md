# DevQE vcenter certificates

- certificates in this directory can expire, but the errors will not be too specific about it and could look like this:
```
k8s.io/client-go/tools/watch/informerwatcher.go:146: Failed to watch *v1.ConfigMap: failed to list *v1.ConfigMap: Get "https://api.rbednar.devqe.ibmc.devcluster.openshift.com:6443/api/v1/namespaces/kube-system/configmaps?fieldSelector=metadata.name%3Dbootstrap&resourceVersion=4554": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster.local, openshift, openshift.default, openshift.default.svc, openshift.default.svc.cluster.local, 172.30.0.1, not api.rbednar.devqe.ibmc.devcluster.openshift.com

```

- they are required for running openshift-install on vsphere
- can be obtained from DevQE vcenter home page under `Download trusted root CA certificates` here: https://vcenter.devqe.ibmc.devcluster.openshift.com/
