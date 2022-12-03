package main

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "github.com/RealHarshThakur/grpc-filestream/protos/filestream"
)

const (
	address = "localhost:50051"
)

func main() {
	// Create a connection to the server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a new client
	c := pb.NewFileStreamServiceClient(conn)

	// Send a request for the file
	filename := "myfile.txt"
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()
	stream, err := c.StreamFile(ctx, &pb.StreamFileRequest{Filename: filename})
	if err != nil {
		log.Fatalf("Failed to stream file: %v", err)
	}

	// Receive the file contents as a stream of chunks
	totalbytes := 0
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to receive chunk: %v", err)
		}
		totalbytes += len(chunk.Chunk)
		log.Print(string(chunk.Chunk))
	}
	log.Printf("Received chunk of %d kilobytes", totalbytes/1024)
}
