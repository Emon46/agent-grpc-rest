package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
	"log"
	"observo/agent-grpc/pb"
)

func newGetRequest() *pb.GetWorkerConfigRequest {
	return &pb.GetWorkerConfigRequest{
		ConfigMeta: &pb.ConfigMeta{
			ConfigMapName:      "obsv-data-plane-config",
			ConfigMapNamespace: "obsv-data-plane",
		},
	}
}

func testGetWorkerConfig(serviceClient pb.ConfigServiceClient) {
	fmt.Println("**************** GET CONFIG **************")
	config, err := serviceClient.GetWorkerConfig(context.TODO(), newGetRequest())
	if err != nil {
		log.Fatal(err)
	}
	jsonConfig, err := json.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonConfig))
}

func newUpdateRequest() *pb.UpdateWorkerConfigRequest {
	transformsConfig := map[string]interface{}{
		"filter_test": map[string]interface{}{
			"condition": "random condition",

			"inputs": []interface{}{"k8s_logs_source"},
			"type":   "filter",
		},
	}

	transformsConfigPB, err := structpb.NewStruct(transformsConfig)
	if err != nil {
		panic(err)
	}
	return &pb.UpdateWorkerConfigRequest{
		ConfigMeta: &pb.ConfigMeta{
			ConfigMapName:      "obsv-data-plane-config",
			ConfigMapNamespace: "obsv-data-plane",
		},
		Config: &pb.Config{
			Transforms: transformsConfigPB,
		},
	}
}

func testUpdateWorkerConfig(serviceClient pb.ConfigServiceClient) {
	fmt.Println("**************** UPDATE CONFIG **************")
	config, err := serviceClient.UpdateWorkerConfig(context.TODO(), newUpdateRequest())
	if err != nil {
		log.Fatal(err)
	}
	jsonConfig, err := json.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonConfig))
}

func main() {
	address := flag.String("address", "", "the server port")
	flag.Parse()
	log.Println("start the server on port", *address)
	log.Printf("dialing %s server address", *address)
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	configServerClient := pb.NewConfigServiceClient(conn)
	//testSearchLaptop(laptopClient)
	testGetWorkerConfig(configServerClient)
	testUpdateWorkerConfig(configServerClient)

}
