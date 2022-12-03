package main

import (
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"

	pb "github.com/RealHarshThakur/grpc-filestream/protos/filestream"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedFileStreamServiceServer
}

func (s *server) StreamFile(req *pb.StreamFileRequest, stream pb.FileStreamService_StreamFileServer) error {
	// Open the file
	file, err := os.Open(req.Filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the file contents
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	// Send the file contents as a stream of chunks
	for len(data) > 0 {
		time.Sleep(5 * time.Second)
		chunkSize := 1 * 1024 // 64KB
		if chunkSize > len(data) {
			chunkSize = len(data)
		}
		chunk := data[:chunkSize]
		if err := stream.Send(&pb.StreamFileResponse{Chunk: chunk}); err != nil {
			return err
		}
		data = data[chunkSize:]
	}

	return nil
}

func main() {
	// Create a listener for incoming connections
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server
	s := grpc.NewServer()

	// Register the FileStreamService with the server
	pb.RegisterFileStreamServiceServer(s, &server{})

	// Start serving requests
	log.Println("Starting server on port", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
