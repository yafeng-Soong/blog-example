package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"resolver-with-watch/proto"
	"resolver-with-watch/resolver"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	connCloseFunc = append(connCloseFunc, closeHelloConn)
}

var (
	helloConn *grpc.ClientConn
	lock      sync.Mutex
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "tester"
	}

	var err error
	defer func() {
		if err != nil {
			log.Println(err.Error())
			w.Write([]byte(err.Error()))
		}
	}()

	sayHelloClient, err := getSayHeloClient()
	if err != nil {
		return
	}

	out, err := sayHelloClient.SayHello(context.Background(), &proto.SayHelloRequest{Name: name})
	if err != nil {
		return
	}

	echo := fmt.Sprintf("%s, now is %s", out.Hello, time.Now().Format("2006-01-02 15:04:05.000"))
	log.Printf("echo content: %s", echo)
	w.Write([]byte(echo))
}

func getSayHeloClient() (proto.HelloServiceClient, error) {
	lock.Lock()
	defer lock.Unlock()

	if helloConn != nil {
		return proto.NewHelloServiceClient(helloConn), nil
	}

	target := fmt.Sprintf("%s://%s", resolver.Scheme, "hello-server")
	var err error
	helloConn, err = grpc.Dial(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		return nil, err
	}

	return proto.NewHelloServiceClient(helloConn), nil
}

func closeHelloConn() {
	if helloConn == nil {
		return
	}

	helloConn.Close()
}
