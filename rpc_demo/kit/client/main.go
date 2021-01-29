package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"micro-go/rpc_demo/pd"
	"time"
)

func main() {
	flag.Parse()

	ctx := context.Background()
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		log.Println("grpc dial err:", err)
	}
	defer conn.Close()

	svc := pd.NewStringServiceClient(conn)

	stringReq := &pd.StringRequest{A: "A", B: "B"}
	reply, _ := svc.Concat(ctx, stringReq)
	fmt.Printf("StringService Concat: %s concat %s =%s", stringReq.A, stringReq.B, reply.Ret)
}
