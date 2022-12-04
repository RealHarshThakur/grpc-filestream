# grpc-kubelog-stream 

If you are using Kubernetes and want to stream logs externally without using the Kubernetes API, this is the tool for you.

It is a gRPC server which accepts a request(Job name, etc) and streams the logs of all pods to the client. A client is also provided to demonstrate how to use the server.


## Build
Server:
```
cd server; go build . ; mv server ../bin
```
Client:
```
cd client; go build . ; mv client ../bin
```

## Usage
* Apply the job manifest to your cluster
```
kubectl apply -f job.yaml
```

* Run the server
```
// explicitly specify the kubeconfig file via env var
export KUBECONFIG=/path/to/kubeconfig
./bin/server
```

* Run the client
```
./bin/client
```

