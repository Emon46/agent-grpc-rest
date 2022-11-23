package grpc

import (
	"k8s.io/client-go/kubernetes"
	"observo/agent-grpc/internal/lib"
	"observo/agent-grpc/pb"
)

type Server struct {
	ServerConfig *lib.ServerConfig
	KubeClient   kubernetes.Interface
	pb.UnimplementedConfigServiceServer
}

// NewServer it will create a new gin api and setup routing for all the api call
func NewServer(config lib.ServerConfig, kubeClient kubernetes.Interface) (*Server, error) {

	server := &Server{
		ServerConfig: &config,
		KubeClient:   kubeClient,
	}

	return server, nil
}
