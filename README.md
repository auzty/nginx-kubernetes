# Nginx Remotehand 

This simple daemon will listen to kubernetes API and write the nginx
configuration

This base images are using nginx https://hub.docker.com/_/nginx/


## Compile 

go build -o k8s-api kubernetes-svc-api.go conf-handler.go

## Run

CLUSTERNAME=cluster.local ./final

### MIsc

- the file will be written to /etc/nginx/conf.d/xxx.conf
- this daemon only create a file, not overwrite it
- you must specify the label to kubernetes service manifest file

example :

```yaml
apiVersion: v1
kind: Service
metadata:
  name: "hello-world"
  namespace: default
  labels:
    app: "hello-world"
    domain: "jajal.io"
    internaldomain: "jajal.dev.io"
spec:
  selector:
    app: "hello-world"
  ports:
    - name: http
      port: 3080
      targetPort: 3001
      protocol: TCP
  sessionAffinity: ClientIP
```

- the internal domain must be set

- the template file **conf.example** are golang template file format 

### Using Kubernetes

the binary will read template file from
**/opt/nginx-k8sapi/conf.template** , to replace it can using mount
volume 

#### Make sure create the service account first

(service account) is an account that created on kubernetes to access the
api inside k8s cluster

**example.yaml**

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: defaultapi
  namespace: default
```

```bash
kubectl create -f example.yaml
```

Check the token name using 

```bash
kubectl describe serviceaccount defaultapi
```

And then add it to nginx manifest file

**nginx.yaml**
```yaml
apiVersion: "extensions/v1beta1"
kind: Deployment
metadata:
  name: nginx
  namespace: default
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx-reload:v2
          imagePullPolicy: IfNotPresent
          volumeMounts:
            - name: apiaccess
              mountPath: /var/run/secrets/kubernetes.io/serviceaccount
              readOnly: true
      volumes:
        - name: apiaccess
          secret:
            defaultMode: 420
            secretName: defaultapi-token-wr531
      restartPolicy: Always
```

# Dockerfile

## Building the images 

```bash
docker build -t nginx-reload .
```

this nginx images will compare the sha1 all of the configuration located
at **/etc/nginx/conf.d/** and if there are some changes, this will
inform nginx to reload the service.
