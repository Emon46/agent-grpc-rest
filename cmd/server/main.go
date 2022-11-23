package main

import (
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"kmodules.xyz/client-go/tools/clientcmd"
	"log"
	"net"
	"observo/agent-grpc/internal/lib"
	grpc_sv "observo/agent-grpc/internal/service/grpc"
	"observo/agent-grpc/internal/service/rest"
	"observo/agent-grpc/pb"
	"sync"
)

func main() {
	config, err := lib.LoadConfig("/configs/")
	if err != nil {
		log.Fatal(err)
	}
	kubeConfig, err := restclient.InClusterConfig()
	if err != nil {
		klog.Fatalln(err)
	}
	clientcmd.Fix(kubeConfig)
	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		klog.Fatalln(err)
	}

	var wg sync.WaitGroup

	grpcServer, err := grpc_sv.NewServer(config, kubeClient)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := runGrpcServer(grpcServer)
		if err != nil {
			log.Fatal("cannot start server: ", err)
		}
	}()

	restServer, err := rest.NewServer(config, kubeClient)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := runRestServer(restServer)
		if err != nil {
			log.Fatal("cannot start server: ", err)
		}
	}()

	wg.Wait()
}

func runGrpcServer(server *grpc_sv.Server) error {
	grpcServer := grpc.NewServer()

	pb.RegisterConfigServiceServer(grpcServer, server)
	listener, err := net.Listen("tcp", server.ServerConfig.GRPCServerAddress)
	if err != nil {
		return err
	}
	log.Println("listening to get and update config server")
	err = grpcServer.Serve(listener)
	if err != nil {
		return err
	}
	return nil
}

func runRestServer(server *rest.Server) error {
	log.Println("starting agent server...")
	err := server.Start(server.ServerConfig.RestServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
	return nil
}
