package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"resolver-with-watch/builder"
	"resolver-with-watch/proto"
	"resolver-with-watch/register"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port         = flag.Int("port", 8080, "")
	registerAddr = flag.String("register", "localhost:2379", "")
	once         sync.Once
	helloConn    *grpc.ClientConn
	helloClient  proto.HelloServiceClient
)

func init() {
	flag.Parse()
	register.InitRegister(*registerAddr)
}

func main() {
	http.HandleFunc("/echo", echo)
	go func() {
		log.Printf("serve at :%d", *port)
		http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	if helloConn != nil {
		helloConn.Close()
	}
}

func echo(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "tester"
	}

	sayHelloClient := getSayHeloClient()
	out, err := sayHelloClient.SayHello(context.Background(), &proto.SayHelloRequest{Name: name})
	if err != nil {
		log.Println(err.Error())
		w.Write([]byte(err.Error()))
		return
	}

	echo := fmt.Sprintf("%s, now is %s", out.Hello, time.Now().Format("2006-01-02 15:04:05.000"))
	log.Printf("echo content: %s", echo)
	w.Write([]byte(echo))
}

func initSayHelloConn() {
	target := fmt.Sprintf("%s://%s", builder.Scheme, "hello-server")
	var err error
	helloConn, err = grpc.Dial(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func getSayHeloClient() proto.HelloServiceClient {
	once.Do(func() {
		initSayHelloConn()
		helloClient = proto.NewHelloServiceClient(helloConn)
	})
	return helloClient
}
