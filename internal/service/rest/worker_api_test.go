package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"io"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"observo/agent-grpc/internal/lib"
	"strings"
	"testing"

	testclient "k8s.io/client-go/kubernetes/fake"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

const configMapStr = `
data_dir: /vector-data-dir
sinks:
    k8s_logs_sink:
        compression: none
        encoding:
            codec: json
        inputs:
            - filter_k8s_logs
        path: /tmp/vector-demo-logs-1-%Y-%m-%d.log
        type: file
    vector_agent_sink:
        address: http://vector-agent:9000
        inputs:
            - internal_log_source
        type: vector
transforms:
    filter_k8s_logs:
        condition: contains(string(.message) ?? "", "no_tag") != true
        inputs:
            - k8s_logs_source
        type: filter
sources:
    internal_log_source:
        type: internal_logs
    k8s_logs_source:
        extra_field_selector: metadata.name==load-test-pod
        type: kubernetes_logs
`

func TestServerGetWorkerConfig(t *testing.T) {
	configMapNamespace := RandomString(4)
	configMapName := RandomString(8)
	clientset := testclient.NewSimpleClientset()
	configMap, err := createDummyConfigMap(configMapName, configMapNamespace, clientset)
	require.NoError(t, err)
	require.NotNil(t, configMap)

	testCases := []struct {
		name          string
		body          gin.H
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "ok",
			body: gin.H{
				"configMapMeta": &ConfigMapMeta{
					ConfigMapName:      configMapName,
					ConfigMapNamespace: configMapNamespace,
				},
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				matchWorkerConfig(t, recorder.Body, configMapStr)
			},
		},
		{
			name: "invalid namespace",
			body: gin.H{
				"configMapMeta": &ConfigMapMeta{
					ConfigMapName:      configMapName,
					ConfigMapNamespace: RandomString(4),
				},
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "invalid config name",
			body: gin.H{
				"configMapMeta": &ConfigMapMeta{
					ConfigMapName:      RandomString(8),
					ConfigMapNamespace: configMapNamespace,
				},
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "missing namespace",
			body: gin.H{
				"configMapMeta": &ConfigMapMeta{
					ConfigMapName: configMapName,
				},
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "missing config name",
			body: gin.H{
				"configMapMeta": &ConfigMapMeta{
					ConfigMapNamespace: configMapNamespace,
				},
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			server := newTestServer(t, clientset)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/config"
			request, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
func matchWorkerConfig(t *testing.T, body *bytes.Buffer, configStr string) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	resultConfig := make(map[string]interface{})
	err = json.Unmarshal(data, &resultConfig)
	require.NoError(t, err)
	require.NotEmpty(t, resultConfig)

	expectedConfig := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(configStr), &expectedConfig)
	require.NoError(t, err)
	require.NotEmpty(t, expectedConfig)
	require.Equal(t, expectedConfig, resultConfig)
}

func TestServerUpdateWorkerConfig(t *testing.T) {
	configMapNamespace := RandomString(4)
	configMapName := RandomString(8)
	clientset := testclient.NewSimpleClientset()
	configMap, err := createDummyConfigMap(configMapName, configMapNamespace, clientset)
	require.NoError(t, err)
	require.NotNil(t, configMap)

	// these are the changes we are going to request to in worker configmap
	transformConfigs := getTransformNewConfigs()

	testCases := []struct {
		name          string
		secret        *core_v1.ConfigMap
		body          gin.H
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "ok",
			body: gin.H{
				"configMapMeta": &ConfigMapMeta{
					ConfigMapName:      configMapName,
					ConfigMapNamespace: configMapNamespace,
				},
				"config": Config{
					Transforms: transformConfigs,
				},
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				checkForNewWorkerConfig(t, recorder.Body, transformConfigs)
			},
		},
		{
			name: "invalid namespace",
			body: gin.H{
				"configMapMeta": &ConfigMapMeta{
					ConfigMapName:      configMapName,
					ConfigMapNamespace: RandomString(4),
				},
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "invalid config name",
			body: gin.H{
				"configMapMeta": &ConfigMapMeta{
					ConfigMapName:      RandomString(8),
					ConfigMapNamespace: configMapNamespace,
				},
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "missing namespace",
			body: gin.H{
				"configMapMeta": &ConfigMapMeta{
					ConfigMapName: configMapName,
				},
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "missing config name",
			body: gin.H{
				"configMapMeta": &ConfigMapMeta{
					ConfigMapNamespace: configMapNamespace,
				},
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			server := newTestServer(t, clientset)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/config"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func getTransformNewConfigs() map[string]interface{} {
	return map[string]interface{}{
		"filter_test": struct {
			Condition string
			Inputs    []string
			Type      string
		}{
			Condition: "random condition",
			Inputs:    []string{"k8s_logs_source"},
			Type:      "filter",
		},
		"filter_test_2": struct {
			Condition string
			Inputs    []string
			Type      string
		}{
			Condition: "random condition",
			Inputs:    []string{"k8s_logs_source"},
			Type:      "filter",
		},
	}
}

func checkForNewWorkerConfig(t *testing.T, body *bytes.Buffer, transformConf map[string]interface{}) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	resultConfig := make(map[string]interface{})
	err = json.Unmarshal(data, &resultConfig)
	require.NoError(t, err)
	require.NotEmpty(t, resultConfig)

	fmt.Println(resultConfig)
	fmt.Println(transformConf)
	transformConfig := resultConfig["transforms"].(map[string]interface{})
	for key, _ := range transformConf {
		require.NotEmpty(t, transformConfig[key])
	}
}

func createDummyConfigMap(name, namespace string, clientset kubernetes.Interface) (*core_v1.ConfigMap, error) {
	config, err := clientset.CoreV1().ConfigMaps(namespace).Create(context.TODO(), getDummyConfigMap(name, namespace), meta_v1.CreateOptions{})
	return config, err
}

func getDummyConfigMap(name, namespace string) *core_v1.ConfigMap {
	data := make(map[string]string)
	data[lib.WorkerConfigFileName] = configMapStr
	return &core_v1.ConfigMap{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}
