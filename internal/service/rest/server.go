package rest

import (
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
	"observo/agent-grpc/internal/lib"
)

type Server struct {
	ServerConfig lib.ServerConfig
	KubeClient   kubernetes.Interface
	router       *gin.Engine
}

// NewServer it will create a new gin api and setup routing for all the api call
func NewServer(config lib.ServerConfig, kubeClient kubernetes.Interface) (*Server, error) {

	server := &Server{
		ServerConfig: config,
		KubeClient:   kubeClient,
	}

	server.setUpRouterWithSubUrl()
	return server, nil
}

func (server *Server) setUpRouterWithSubUrl() {
	router := gin.Default()
	router.GET("/config", server.GetWorkerConfig)
	router.POST("/config", server.UpdateWorkerConfig)

	server.router = router
}
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
