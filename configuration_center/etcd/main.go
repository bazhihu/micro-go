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
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Printf("connect to etcd failed, err: %v\n", err)
		return
	}

	fmt.Println("connect to etcd success")
	defer cli.Close()

	//put
	_ = PutValue("/ddd", `{"ddd":1}`)

	// get
	val := GetValue("/ddd")
	log.Printf("value: %v\n", val)

	// watch this key
	c := make(chan []byte, 1)
	go Watch("/ddd", c)

	go func() {
		for {
			log.Println(string(<-c))
		}
	}()

	_ = PutValue("/ddd", `{"ddd":2}`)

	// lease租约
	// 创建一个5秒租约
	resp, err := EtcdClient.Grant(context.TODO(), 5)
	if err != nil {
		log.Fatal("Grant-err:", err)
	}

	// 5秒之后，key就会被移除
	_, err = EtcdClient.Put(context.TODO(), "/lmh/", "lmh", clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatalf("put-grant-err:", err)
	}

	//log.Fatal("time:", string(GetValue("/lmh/")))

	time.Sleep(6 * time.Second)
	log.Fatal("6 second:", string(GetValue("/lmh/")))
}

var EtcdClient *clientv3.Client

func init() {
	var err error
	EtcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("connect to etcd failed, err: %v\n", err)
		panic(err)
	}
}

// 设置键值
func PutValue(key, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	_, err := EtcdClient.Put(ctx, key, value)
	cancel()
	if err != nil {
		log.Fatalf("etcd put err:%v\n", err)
	}
	return err
}

// 获取键值
func GetValue(key string) (value []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	resp, err := EtcdClient.Get(ctx, key)
	cancel()
	if err != nil {
		log.Fatalf("etcd get err:%v\n", err)
		return
	}
	for _, ev := range resp.Kvs {
		value = ev.Value
		break
	}
	return
}

// watch 键值 通过c 传递
func Watch(key string, c chan<- []byte) {
	watch := EtcdClient.Watch(context.Background(), key)
	for wresp := range watch {
		for _, v := range wresp.Events {
			c <- v.Kv.Value
			log.Printf("watch key:%s change value %v", key, string(v.Kv.Value))
		}
	}
}
