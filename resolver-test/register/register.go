package register

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var client *clientv3.Client

func InitRegister(addr string) {
	var err error
	client, err = clientv3.NewFromURL(addr)
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

func Register(serviceName string, addr string) (string, error) {
	key := fmt.Sprintf("%s/%v", serviceName, time.Now().UnixNano())
	_, err := client.Put(context.Background(), key, addr)
	return key, err
}

func DeRegister(instanceName string) {
	client.Delete(context.Background(), instanceName)
}
