package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Printf("connect to etcd failed, err: %v\n", err)
		return
	}

	fmt.Println("connect to etcd success")
	defer cli.Close()

	//put
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, "ddd", `{"ddd":1}`)
	cancel()
	if err != nil {
		log.Printf("put to etcd failed, err:%v\n", err)
		return
	}

	// get
	ctx, cancel = context.WithTimeout(ctx, time.Second)
	resp, err := cli.Get(ctx, "ddd")
	cancel()
	if err != nil {
		log.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	for _, ev := range resp.Kvs {
		log.Printf("%s:%s\n", ev.Key, ev.Value)
	}
}
