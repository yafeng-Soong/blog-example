package register

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

const dialTimeout = 5 * time.Second

var (
	client       *clientv3.Client
	instanceName string
)

func InitRegister(addr string) {
	var err error
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: dialTimeout,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	if err != nil {
		panic(fmt.Sprintf("init etcd client error: %s", err.Error()))
	}
}

func QueryAddress(serviceName string) (addrs []string) {
	resp, err := client.Get(context.Background(), serviceName, clientv3.WithPrefix())
	if err != nil {
		return
	}

	for _, kv := range resp.Kvs {
		addrs = append(addrs, string(kv.Value))
	}
	log.Println(fmt.Sprintf("found service %s in: ", serviceName), addrs)
	return
}

func Register(serviceName string, addr string) error {
	instanceName = fmt.Sprintf("%s/%v", serviceName, time.Now().UnixNano())
	_, err := client.Put(context.Background(), instanceName, addr)
	return err
}

func DeRegister() {
	if instanceName == "" {
		return
	}

	client.Delete(context.Background(), instanceName)
}

func CloseRegister() {
	client.Close()
}
