package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	con "github.com/coreos/etcd/clientv3/concurrency"
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
	//Grant("/lmh/", "lmh", 5)

	// keepalive
	//KeepAlive("/lll/", "lll", 10)

	var em1 EtcdMutex // 分布式锁1
	var em2 EtcdMutex // 分布式锁2

	// 枷锁1
	em1.Lock("/lock_demo_1/")
	log.Println("lock one key /lock_demo_1/")
	go func() { // 10秒后释放锁
		defer em1.UnLock()
		time.Sleep(10 * time.Second)
	}()

	if err := em2.Lock("/lock_demo_1/"); err != nil {
		log.Fatal("lock two err:", err)
	}
	defer em2.UnLock()

	log.Println("lock two key /lock_demo_1/")
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

// 创建租约
func Grant(key, value string, ttl int64) {
	resp, err := EtcdClient.Grant(context.TODO(), ttl)
	if err != nil {
		log.Fatalf("etcd grant err: %v\n", err)
	}
	_, err = EtcdClient.Put(context.TODO(), key, value, clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatalf("etcd grant put err: %v\n", err)
	}

	// 如果读取了 租约就消失了
	//GetValue(key)
}

// 长链接
func KeepAlive(key, value string, ttl int64) {
	resp, err := EtcdClient.Grant(context.TODO(), ttl)
	if err != nil {
		log.Fatalf("etcd grant err: %v\n", err)
	}
	_, err = EtcdClient.Put(context.TODO(), key, value, clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatalf("etcd grant put err: %v\n", err)
	}

	// the key will by kept forever
	ch, kaerr := EtcdClient.KeepAlive(context.TODO(), resp.ID)
	if kaerr != nil {
		log.Fatalf("etcd grant KeepAlive err: %v\n", kaerr)
	}

	var i int
	for {
		i++
		if i > 10 {
			break
		}
		ka := <-ch
		log.Printf("ttl: %d\n", ka.TTL)
	}
}

type EtcdMutex struct {
	session *con.Session
	mutex   *con.Mutex
}

// 分布式枷锁
func (em *EtcdMutex) Lock(key string) error {
	var (
		err error
	)

	if em.session, err = con.NewSession(EtcdClient); err != nil {
		log.Printf("Create NewSession err %v\n", err)
	}

	em.mutex = con.NewMutex(em.session, key)
	if err = em.mutex.Lock(context.TODO()); err != nil {
		log.Fatalf("lock fatal %v\n", err)
	}

	return err
}

// 分布式解锁
func (em *EtcdMutex) UnLock() error {
	defer em.session.Close()
	var (
		err error
	)
	if err = em.mutex.Unlock(context.TODO()); err != nil {
		log.Fatalf("unlock fatal %v\n", err)
	}
	return err
}
