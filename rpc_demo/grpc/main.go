package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"micro-go/rpc_demo/grpc/service"
	"micro-go/rpc_demo/pd"
	"net"
)

func main() {
	var (
		port = flag.Int("service.port", 1234, "service port")
	)

	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	stringService := new(service.StringService)
	pd.RegisterStringServiceServer(grpcServer, stringService)
	grpcServer.Serve(lis)
}
