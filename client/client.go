package main

import (
	"context"
	"io"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	pb "github.com/RealHarshThakur/grpc-kubelog-stream/protos/stream"
	"github.com/sirupsen/logrus"
)

const (
	address = "localhost:50051"
)

func main() {

	log := SetupLogging()
	// Create a connection to the server
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:    30 * time.Second, // client ping server if no activity for this long
		Timeout: 20 * time.Second,
	}))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a new client
	c := pb.NewJobLogsServiceClient(conn)

	// Send a request for the file
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()
	stream, err := c.GetJobLogs(ctx, &pb.GetJobLogsRequest{Name: "pi-with-ttl", Namespace: "default"})
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
		totalbytes += len(chunk.GetLogs())
		log.Print(string(chunk.GetLogs()))
	}
	log.Printf("Received chunk of %d kilobytes", totalbytes/1024)
}

// SetupLogging sets up the logging for the router daemon
func SetupLogging() *logrus.Logger {
	// Logging create logging object
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)
	// log.SetReportCaller(true)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		// CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
		// 	fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
		// 	return "", fileName
		// },
	})

	return log
}
