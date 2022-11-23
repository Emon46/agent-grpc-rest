package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes"
	"observo/agent-grpc/internal/lib"
	"os"
	"testing"
)

func newTestServer(t *testing.T, clientset kubernetes.Interface) *Server {
	config := lib.ServerConfig{
		RestServerAddress: "0.0.0.0:8080",
	}
	server, err := NewServer(config, clientset)
	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
