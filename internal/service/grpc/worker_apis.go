package grpc

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/structpb"
	"gopkg.in/yaml.v3"
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

	workerConfig, err := lib.UpsertWorkerConfig(configMap, &lib.WorkerConfig{
		Sources:    req.Config.Sources.AsMap(),
		Transforms: req.Config.Transforms.AsMap(),
		Sinks:      req.Config.Sinks.AsMap(),
	})
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
