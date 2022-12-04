package main

import (
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	pb "github.com/RealHarshThakur/grpc-kubelog-stream/protos/stream"
)

const (
	port = ":50051"
)

var kubeconfig string

func main() {
	kubeconfig = os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		log.Fatal("KUBECONFIG environment variable not set")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Create a listener for incoming connections
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server
	s := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute, // <--- This fixes it!
			Time:              5 * time.Minute,
			Timeout:           5 * time.Minute,
		}),
	)

	// Register the FileStreamService with the server
	pb.RegisterFileStreamServiceServer(s, &fileServer{})

	pb.RegisterJobLogsServiceServer(s, &jobLogsServer{
		KClient: clientset,
	})

	// Start serving requests
	log.Println("Starting server on port", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
