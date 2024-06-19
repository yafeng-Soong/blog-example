package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"resolver-test/builder"
	"resolver-test/proto"
	"resolver-test/register"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	name         = flag.String("name", "tester", "")
	registerAddr = flag.String("register", "localhost:2379", "")
)

func init() {
	flag.Parse()
}

// go run client.go [-name echo_name] [-register 127.0.0.1:2379]
func main() {
	register.InitRegister(*registerAddr)
	defer register.CloseRegister()

	target := fmt.Sprintf("%s://%s", builder.Scheme, "hello-server")
	conn, err := grpc.Dial(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer conn.Close()

	client := proto.NewHelloServiceClient(conn)
	for i := 0; i < 10; i++ {
		out, err := client.SayHello(context.Background(), &proto.SayHelloRequest{Name: *name})
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Println("echo: ", out.Hello)
		time.Sleep(2 * time.Second)
	}
}
