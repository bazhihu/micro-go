package main

import (
	"fmt"
	"log"
	"micro-go/rpc_demo/service"
	"net/rpc"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	stringReq := &service.StringRequest{"A", "B"}

	// synchronous call
	var reply string
	err = client.Call("StringService.Concat", stringReq, &reply)
	if err != nil {
		log.Fatal("Concat error:", err)
	}

	fmt.Printf("stringService Concat: %s concat %s = %s\n", stringReq.A, stringReq.B, reply)

	// go
	stringReq = &service.StringRequest{"ACD", "BDF"}
	call := client.Go("StringService.Diff", stringReq, &reply, nil)
	_ = <-call.Done
	fmt.Printf("stringService Diff: %s diff %s = %s\n", stringReq.A, stringReq.B, reply)
}
