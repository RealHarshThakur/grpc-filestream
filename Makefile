protos: 
  protoc -I . filestream.proto --go-grpc_out=.  --go_out=.

server: 
   cd server; go build . ; mv server ../bin

# client:
#   go build -o bin/client client
