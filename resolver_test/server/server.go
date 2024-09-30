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

var (
	registerAddr = flag.String("register", "localhost:2379", "")
	localhost    string
)

func init() {
	flag.Parse()
	localhost = getLocalIP()
}

// go run server.go [-register 127.0.0.1:2379]
func main() {
	register.InitRegister(*registerAddr)
	defer register.CloseRegister()

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:0", localhost))
	if err != nil {
		log.Fatal(err.Error())
	}

	s := grpc.NewServer()
	addr := listener.Addr().String()
	proto.RegisterHelloServiceServer(s, &helloServer{addr: addr})
	log.Println("hello-server serve at: ", addr)

	if err := register.Register("hello-server", addr); err != nil {
		log.Fatal(err.Error())
	}

	defer register.DeRegister()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := s.Serve(listener); err != nil {
			log.Println(err.Error())
			close(signalChan)
		}
	}()

	<-signalChan
	s.GracefulStop()
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
	proto.UnimplementedHelloServiceServer
	addr string
}

func (s *helloServer) SayHello(ctx context.Context, req *proto.SayHelloRequest) (*proto.SayHelloResponse, error) {
	log.Println("serve ", req.Name)
	return &proto.SayHelloResponse{
		Hello: fmt.Sprintf("hello %s, from server %s", req.Name, s.addr),
	}, nil
}
