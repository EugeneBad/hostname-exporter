## Hostname Exporter

A light-weight golang web service that exports the hostname of the host it's running on; which can be kubernetes nodes or vanilla virtual machines.

Additionally, it exposes prometheus metrics on the `/metrics` endpoint; to track request count, both failed and successful.

### Folder structure
```
.
├── Dockerfile
├── cmd/
│   ├── go.mod
│   ├── go.sum
|   ├── main.go
|   ├── server.go
│   └── server_test.go
├── kind/
│   └── kind-cluster.yml
└── kubernetes/
    ├── ingress.yml
    ├── namespace.yml
    ├── deployment.yml
    └── service.yml
```
### Endpoints
* `/hostname`: returns the hostname and timestamp as json
    Sample json response:
    ```json
    {"timestamp":"2023-02-03T15:37:03Z","hostname":"kind-control-plane"}
    ```
* `/metrics`: returns prometheus metrics
    ```text
    # HELP app_request_count_total The total number of requests processed
    # TYPE app_request_count_total counter
    app_request_count_total{code="200",method="GET"} 3
    app_request_count_total{code="400",method="POST"} 1
    ```
### Setup
Install golang dependencies using the `go mod` utility:
```
$ cd cmd
$ go mod tidy
```

```
$ go get github.com/prometheus/client_golang/prometheus/promauto
$ go get github.com/prometheus/client_golang/prometheus
$ go get github.com/prometheus/client_golang/prometheus/promhttp
```
### Unit-testing
The server's handler function is unit-tested using the golang test package: `testify`
```
$ cd cmd
$ go test

PASS
ok      adjust  0.004s
```

### Building
The `Dockerfile` packages the service into an alpine based docker image. This keeps the image size small and quick to deploy:

```
docker build -t eugenebad/hostname-exporter:0.0.1c .
```
Push image:
```
docker push eugenebad/hostname-exporter:0.0.1c
```

> Inorder to push to a different docker registry, rename the image accordingly:
> ```docker build -t <registry_name>/hostname-exporter:0.0.1c .```


### Running locally
The service can be tested locally by running the docker image with a port mapped to the listening address `9090`
```
docker run -it -p "9090:9090" eugenebad/hostname-exporter:0.0.1c
```

This exposes the service on the local loopback address `127.0.0.1` at port `9090`

```
$ curl 127.0.0.1:9090/hostname
```

### Running on kubernetes
The service can be deployed on a kubernetes cluster using the manifest files in `/kubernetes`

Such a test environment can be bootstrapped locally using the [`kind`](https://kind.sigs.k8s.io/) utility.

##### 1. Installation (on Linux)
Download and install the kind binary
```
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.17.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind
```

##### 2. Create kind cluster
From repo directory:
```
$ kind create cluster --config kind/kind-cluster.yml
```
> Note: This sets your current kubectl context to be `kind-kind`

##### 3. Install Nginx ingress controller
Inorder to route external requests to services running inside the cluster:
```
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
```
##### 4. Deploy the hostname-exporter service
1. Create namespace:
```
$ kubectl create -f kubernetes/namespace.yml
```
Before deploying the application, note that the `kubernetes/ingress.yml` file has a **host** attribute which defines the domain name on which the service will be accessed. You can set this to your liking, so long as it resolves to the local address `127.0.0.1`. Which can be achieved by using the `/etc/hosts` file.

> Ideally, for non-local deployment, dns resolution should be handled by a dedicated DNS server.

2. Deploy application:
```
$ kubectl apply -f kubernetes/
```

This exposes the service on `hostname-exporter.local` at port `80`

```
$ curl hostname-exporter.local/hostname
```

##### 5. Cleanup

```
kind delete cluster
```