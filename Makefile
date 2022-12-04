protos: 
  protoc -I . filestream.proto --go-grpc_out=.  --go_out=.

# server: 
#   go build -o bin/server server

# client:
#   go build -o bin/client client
