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

func main() {
	name := flag.String("name", "tester", "")
	registerAddr := flag.String("register", "localhost:2379", "")
	flag.Parse()

	target := fmt.Sprintf("%s://%s", builder.Scheme, "hello-server")
	register.InitRegister(*registerAddr)
	conn, err := grpc.Dial(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer conn.Close()

	client := proto.NewHelloClient(conn)
	for i := 0; i < 10; i++ {
		out, err := client.SayHello(context.Background(), &proto.SayHelloRequest{Name: *name})
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Println("echo: ", out.Echo)
		time.Sleep(2 * time.Second)
	}
}
