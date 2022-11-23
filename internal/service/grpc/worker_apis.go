package grpc

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/structpb"
	"gopkg.in/yaml.v3"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"observo/agent-grpc/internal/lib"
	"observo/agent-grpc/pb"
)

func (s *Server) UpdateWorkerConfig(ctx context.Context, req *pb.UpdateWorkerConfigRequest) (*pb.WorkerConfigResponse, error) {
	if req.ConfigMeta == nil {
		return nil, fmt.Errorf("ConfigMeta spec can't be nill in request struct")
	}

	configMap, err := s.KubeClient.CoreV1().ConfigMaps(req.ConfigMeta.GetConfigMapNamespace()).Get(ctx, req.ConfigMeta.GetConfigMapName(), meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	workerConfig, err := UpsertWorkerConfig(configMap, req.Config)
	if err != nil {
		return nil, err
	}

	err = lib.UpdateConfigMap(ctx, s.KubeClient, configMap, workerConfig)
	if err != nil {
		return nil, err
	}

	configPb, err := structpb.NewStruct(workerConfig)
	if err != nil {
		return nil, err
	}
	return &pb.WorkerConfigResponse{
		Config: configPb,
	}, nil
}

func (s *Server) GetWorkerConfig(ctx context.Context, req *pb.GetWorkerConfigRequest) (*pb.WorkerConfigResponse, error) {
	if req.ConfigMeta == nil {
		return nil, fmt.Errorf("ConfigMeta spec can't be nill in request struct")
	}
	configMap, err := s.KubeClient.CoreV1().ConfigMaps(req.ConfigMeta.GetConfigMapName()).Get(ctx, req.ConfigMeta.GetConfigMapName(), meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	workerConfig := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(configMap.Data[lib.WorkerConfigFileName]), &workerConfig)
	if err != nil {
		return nil, err
	}

	configPb, err := structpb.NewStruct(workerConfig)
	if err != nil {
		return nil, err
	}
	return &pb.WorkerConfigResponse{
		Config: configPb,
	}, nil
}

func UpsertWorkerConfig(configMap *core_v1.ConfigMap, reqConfig *pb.Config) (map[string]interface{}, error) {
	workerConfig := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(configMap.Data[lib.WorkerConfigFileName]), &workerConfig)
	if err != nil {
		return nil, err
	}

	// add new source configs
	if reqConfig.Sources != nil {
		workerConfig[lib.WorkerConfigSourcesStr] = lib.UpsertTypeSpecificWorkerConfig(workerConfig[lib.WorkerConfigSourcesStr].(map[string]interface{}),
			reqConfig.Sources.AsMap())
	}

	// add new transforms configs
	if reqConfig.Transforms != nil {
		workerConfig[lib.WorkerConfigTransformsStr] = lib.UpsertTypeSpecificWorkerConfig(workerConfig[lib.WorkerConfigTransformsStr].(map[string]interface{}),
			reqConfig.Transforms.AsMap())
	}

	// add new sinks configs
	if reqConfig.Sinks != nil {
		workerConfig[lib.WorkerConfigSinksStr] = lib.UpsertTypeSpecificWorkerConfig(workerConfig[lib.WorkerConfigSinksStr].(map[string]interface{}),
			reqConfig.Sinks.AsMap())
	}

	return workerConfig, nil
}
