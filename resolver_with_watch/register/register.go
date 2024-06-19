package register

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

const (
	dialTimeout = 5 * time.Second
	logPrefix   = "[register] "
)

var client *clientv3.Client

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

func QueryAddress(serviceName string) map[string]resolver.Address {
	resp, err := client.Get(context.Background(), serviceName, clientv3.WithPrefix())
	if err != nil {
		return nil
	}

	res := make(map[string]resolver.Address)
	for _, kv := range resp.Kvs {
		res[string(kv.Key)] = resolver.Address{Addr: string(kv.Value)}
	}
	return res
}

func WatchAddress(ctx context.Context,
	cc resolver.ClientConn, serviceName string,
	wacthCallBack func(resolver.ClientConn, string, *clientv3.Event),
) <-chan bool {
	over := make(chan bool)
	go func() {
		for {
			wch := client.Watch(ctx, serviceName, clientv3.WithPrefix())
			select {
			case <-ctx.Done():
				log.Println(logPrefix, "watch over")
				over <- true
				return
			default:
			}

			for wres := range wch {
				if err := wres.Err(); err != nil {
					log.Println(logPrefix, "wacth error: ", err.Error())
					continue
				}

				for _, ev := range wres.Events {
					log.Printf("%swatch event, Type: %s, Key: %s, Value: %s", logPrefix, ev.Type, ev.Kv.Key, ev.Kv.Value)
					wacthCallBack(cc, serviceName, ev)
				}
			}
		}
	}()
	return over
}

func Register(serviceName string, addr string) (string, error) {
	key := fmt.Sprintf("%s/%v", serviceName, time.Now().UnixNano())
	_, err := client.Put(context.Background(), key, addr)
	return key, err
}

func DeRegister(instanceName string) {
	client.Delete(context.Background(), instanceName)
}

func CloseRegister() {
	client.Close()
}
