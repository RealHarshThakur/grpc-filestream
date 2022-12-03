protos:
  protoc -I . filestream.proto --go_out=plugins=grpc:.
