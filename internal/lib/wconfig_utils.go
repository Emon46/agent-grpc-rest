package lib

import (
	"context"
	"gopkg.in/yaml.v3"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	core_util "kmodules.xyz/client-go/core/v1"
)

type WorkerConfig struct {
	Sources    map[string]interface{} `json:"sources"`
	Transforms map[string]interface{} `json:"transforms"`
	Sinks      map[string]interface{} `json:"sinks"`
}

func UpsertWorkerConfig(configMap *core_v1.ConfigMap, reqConfig *WorkerConfig) (map[string]interface{}, error) {
	workerConfig := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(configMap.Data[WorkerConfigFileName]), &workerConfig)
	if err != nil {
		return nil, err
	}

	// add new source configs
	if reqConfig.Sources != nil {
		workerConfig[WorkerConfigSourcesStr] = UpsertTypeSpecificWorkerConfig(workerConfig[WorkerConfigSourcesStr].(map[string]interface{}),
			reqConfig.Sources)
	}

	// add new transforms configs
	if reqConfig.Transforms != nil {
		workerConfig[WorkerConfigTransformsStr] = UpsertTypeSpecificWorkerConfig(workerConfig[WorkerConfigTransformsStr].(map[string]interface{}),
			reqConfig.Transforms)
	}

	// add new sinks configs
	if reqConfig.Sinks != nil {
		workerConfig[WorkerConfigSinksStr] = UpsertTypeSpecificWorkerConfig(workerConfig[WorkerConfigSinksStr].(map[string]interface{}),
			reqConfig.Sinks)
	}

	return workerConfig, nil
}

func UpsertTypeSpecificWorkerConfig(typeSpecificConfig map[string]interface{}, typeSpecificReqConfig map[string]interface{}) map[string]interface{} {
	if typeSpecificConfig == nil {
		typeSpecificConfig = make(map[string]interface{})
	}
	for key, data := range typeSpecificReqConfig {
		typeSpecificConfig[key] = data
	}
	return typeSpecificConfig
}

func UpdateConfigMap(ctx context.Context, kubeClient kubernetes.Interface, configMap *core_v1.ConfigMap, workerConfig map[string]interface{}) error {
	workerDataYaml, err := yaml.Marshal(&workerConfig)
	if err != nil {
		return err
	}
	// now patch the configmap with new updated worker config
	_, _, err = core_util.PatchConfigMap(ctx, kubeClient, configMap, func(in *core_v1.ConfigMap) *core_v1.ConfigMap {
		in.Data[WorkerConfigFileName] = string(workerDataYaml)
		return in
	}, meta_v1.PatchOptions{})
	if err != nil {
		return err
	}
	return nil
}
