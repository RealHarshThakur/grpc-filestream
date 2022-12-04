package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	pb "github.com/RealHarshThakur/grpc-kubelog-stream/protos/stream"
)

type jobLogsServer struct {
	pb.UnimplementedJobLogsServiceServer
	Log     log.Logger
	KClient *kubernetes.Clientset
}

func (d *jobLogsServer) GetJobLogs(req *pb.GetJobLogsRequest, res pb.JobLogsService_GetJobLogsServer) error {
	pods, err := d.KClient.CoreV1().Pods(req.Namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: "job-name=" + req.Name,
	})
	if err != nil {
		d.Log.Printf("error listing pods: %v", err)
		return err
	}

	// Use a waitgroup to wait for all pod logs to be streamed
	var wg sync.WaitGroup
	for _, pod := range pods.Items {
		wg.Add(1)
		go func(name string) {
			err := streamPodLogs(d.KClient, d.Log, &wg, res, name, req.Namespace)
			if err != nil {
				err = fmt.Errorf("error streaming pod logs: %v", err)
				fmt.Println(err)
			}
		}(pod.Name)
	}
	wg.Wait()

	return nil
}

func streamPodLogs(client *kubernetes.Clientset, log log.Logger, wg *sync.WaitGroup, res pb.JobLogsService_GetJobLogsServer, podName, namespace string) error {
	podLogs, err := getPodStream(client, podName, namespace)
	if err != nil {
		log.Printf("error getting pod logstream: %v", err)
		return err
	}
	defer podLogs.Close()

	for {
		// Either exit when the context is done or when the pod stream is closed indicating that all logs have been streamed
		buf := make([]byte, 1024)
		n, err := podLogs.Read(buf)
		if err != nil {
			if res.Context().Err() != nil {
				break
			}
			if err == io.EOF {
				break
			}
			log.Printf("error reading pod log stream: %v", err)
			return err
		}
		if n == 0 {
			// sleep to wait for incoming logs
			time.Sleep(5 * time.Second)
			continue
		}
		msg := fmt.Sprintf("[%s]: %s", podName, string(buf))
		if err := res.Send(&pb.GetJobLogsResponse{Logs: msg}); err != nil {
			err = fmt.Errorf("error sending logs: %v", err)
			log.Printf("error sending logs: %v", err)
			return err
		}
	}
	wg.Done()
	return nil
}

func getPodStream(client *kubernetes.Clientset, podName string, namespace string) (io.ReadCloser, error) {
	logReq := client.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{Follow: true})
	podLogs, err := logReq.Stream(context.Background())
	if err != nil {
		return nil, err
	}
	return podLogs, nil
}
