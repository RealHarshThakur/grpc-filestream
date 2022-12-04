package main

import (
	"io/ioutil"
	"os"
	"time"

	pb "github.com/RealHarshThakur/grpc-kubelog-stream/protos/stream"
)

type fileServer struct {
	pb.UnimplementedFileStreamServiceServer
}

func (s *fileServer) StreamFile(req *pb.StreamFileRequest, stream pb.FileStreamService_StreamFileServer) error {
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
