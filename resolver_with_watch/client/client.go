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
	"resolver-with-watch/proto/protoconnect"
	"resolver-with-watch/register"
	"sync"
	"syscall"
	"time"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port         *int
	registerAddr *string
	once         sync.Once
	helloConn    *grpc.ClientConn
	helloClient  proto.HelloServiceClient
)

func main() {
	port = flag.Int("port", 8080, "")
	registerAddr = flag.String("register", "localhost:2379", "")
	flag.Parse()

	mux := http.NewServeMux()
	path, handler := protoconnect.NewEchoServiceHandler(&echoServer{})
	mux.Handle(path, handler)
	go func() {
		log.Printf("serve at :%d", *port)
		http.ListenAndServe(
			fmt.Sprintf(":%d", *port),
			h2c.NewHandler(mux, &http2.Server{}),
		)
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	if helloConn != nil {
		helloConn.Close()
	}
}

type echoServer struct {
	protoconnect.UnimplementedEchoServiceHandler
}

func (s *echoServer) Echo(
	ctx context.Context,
	req *connect.Request[proto.EchoRequest],
) (*connect.Response[proto.EchoResponse], error) {
	client := getSayHeloClient()
	out, err := client.SayHello(context.Background(), &proto.SayHelloRequest{Name: req.Msg.Name})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	echo := fmt.Sprintf("%s, now is %s", out.Hello, time.Now().Format("2006-01-02 15:04:05.000"))
	log.Printf("echo content: %s", echo)
	return connect.NewResponse(&proto.EchoResponse{Echo: echo}), nil
}

func initSayHelloConn() {
	target := fmt.Sprintf("%s://%s", builder.Scheme, "hello-server")
	register.InitRegister(*registerAddr)
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
