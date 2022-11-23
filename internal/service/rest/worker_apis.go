package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"observo/agent-grpc/internal/lib"
)

type ConfigMapMeta struct {
	ConfigMapName      string `json:"configMapName" binding:"required"`
	ConfigMapNamespace string `json:"configMapNamespace" binding:"required"`
}
type Config struct {
	Sources    map[string]interface{} `json:"sources"`
	Transforms map[string]interface{} `json:"transforms"`
	Sinks      map[string]interface{} `json:"sinks"`
}

type getWorkerConfigRequest struct {
	ConfigMapMeta ConfigMapMeta `json:"configMapMeta" binding:"required"`
}

func (server *Server) GetWorkerConfig(ctx *gin.Context) {
	var req getWorkerConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	configMap, err := server.KubeClient.CoreV1().ConfigMaps(req.ConfigMapMeta.ConfigMapNamespace).Get(context.TODO(), req.ConfigMapMeta.ConfigMapName, meta_v1.GetOptions{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	workerConfig := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(configMap.Data[lib.WorkerConfigFileName]), &workerConfig)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, workerConfig)

}

type updateWorkerConfigRequest struct {
	ConfigMapMeta ConfigMapMeta `json:"configMapMeta" binding:"required"`
	Config        Config        `json:"config"`
}

func (server *Server) UpdateWorkerConfig(ctx *gin.Context) {
	var req updateWorkerConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get the current worker config by getting the config map
	configMap, err := server.KubeClient.CoreV1().ConfigMaps(req.ConfigMapMeta.ConfigMapNamespace).Get(context.TODO(), req.ConfigMapMeta.ConfigMapName, meta_v1.GetOptions{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	workerConfig, err := UpsertWorkerConfig(configMap, req.Config)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = lib.UpdateConfigMap(ctx, server.KubeClient, configMap, workerConfig)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, workerConfig)
}

func UpsertWorkerConfig(configMap *core_v1.ConfigMap, reqConfig Config) (map[string]interface{}, error) {
	workerConfig := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(configMap.Data[lib.WorkerConfigFileName]), &workerConfig)
	if err != nil {
		return nil, err
	}

	// add new source configs
	if reqConfig.Sources != nil {
		workerConfig[lib.WorkerConfigSourcesStr] = lib.UpsertTypeSpecificWorkerConfig(workerConfig[lib.WorkerConfigSourcesStr].(map[string]interface{}),
			reqConfig.Sources)
	}

	// add new transforms configs
	if reqConfig.Transforms != nil {
		workerConfig[lib.WorkerConfigTransformsStr] = lib.UpsertTypeSpecificWorkerConfig(workerConfig[lib.WorkerConfigTransformsStr].(map[string]interface{}),
			reqConfig.Transforms)
	}

	// add new sinks configs
	if reqConfig.Sinks != nil {
		workerConfig[lib.WorkerConfigSinksStr] = lib.UpsertTypeSpecificWorkerConfig(workerConfig[lib.WorkerConfigSinksStr].(map[string]interface{}),
			reqConfig.Sinks)
	}

	return workerConfig, nil
}
