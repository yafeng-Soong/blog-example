package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"resolver-test/proto"
	"resolver-test/register"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	registerAddr := flag.String("register", "localhost:2379", "")
	flag.Parse()
	ip := getLocalIP()
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:0", ip))
	if err != nil {
		log.Fatal(err.Error())
	}

	s := grpc.NewServer()
	proto.RegisterHelloServer(s, &helloServer{})
	addr := listener.Addr().String()
	log.Println("hello-server serve at: ", addr)
	register.InitRegister(*registerAddr)
	ins, err := register.Register("hello-server", addr)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer register.DeRegister(ins)

	signalChan := make(chan os.Signal, 1)
	go func() {
		err = s.Serve(listener)
		if err != nil {
			log.Fatal(err.Error())
		}
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal("get net interface address failed, err = ", err.Error())
	}
	for _, addr := range addrs {
		if ip, ok := addr.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				return ip.IP.String()
			}
		}
	}
	return ""
}

type helloServer struct {
	proto.UnimplementedHelloServer
}

func (s *helloServer) SayHello(ctx context.Context, req *proto.SayHelloRequest) (*proto.SayHelloResponse, error) {
	log.Println("serve ", req.Name)
	return &proto.SayHelloResponse{
		Echo: "hello " + req.Name,
	}, nil
}
