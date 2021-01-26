package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"micro-go/rpc_demo/pd"
)

func main() {
	serviceAddress := "127.0.0.1:1234"
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())
	if err != nil {
		panic("connect error")
	}
	defer conn.Close()
	stringClient := pd.NewStringServiceClient(conn)
	stringReq := &pd.StringRequest{A: "A", B: "B"}
	reply, _ := stringClient.Concat(context.Background(), stringReq)
	fmt.Printf("StringService Concat: %s concat %s =%s", stringReq.A, stringReq.B, reply.Ret)
}
