# vstorage-flexvol
K8s flexvolume driver for vstorage

## to install

In order to use the flexvolume driver, you'll need to install it on every node you want to use ploop on in the kubelet `volume-plugin-dir`. By default this is `/usr/libexec/kubernetes/kubelet-plugins/volume/exec/`

You need a directory for the volume driver vendor, so create it:

```
mkdir -p /usr/libexec/kubernetes/kubelet-plugins/volume/exec/virtuozzo~vstorage-fv
```

Then drop the binary in there:

```
mv ploop /usr/libexec/kubernetes/kubelet-plugins/volume/exec/virtuozzo~vstorage-fv/vstorage-fv
```

### Pod Config

An example pod config would look like this:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx-ploop
spec:
  containers:
  - name: nginx
    image: nginx
    volumeMounts:
    - name: test
      mountPath: /tmp/mnt
    ports:
    - containerPort: 80
  volumes:
  - name: test
    flexVolume:
      driver: "virtuozzo/vstorage-fv" # this must match your vendor dir
      options:
        clusterName: stor1
        clusterPassword: passw0rd
        hostMountPath: "/tmp/mnt"
```


