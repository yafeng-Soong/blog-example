package main

import (
	"flag"
	"log"
	"net/http"
	"tiny-gateway/register"
	"tiny-gateway/rpc"
)

func init() {
	flag.Parse()
}

var registerAddr = flag.String("register", "localhost:2379", "")

func main() {
	register.InitRegister(*registerAddr)
	defer register.CloseRegister()

	http.HandleFunc("/", rpc.HandleRPC)

	log.Println("gRPC gateway is running on port 8080")
	log.Println("Example request: curl -X POST http://localhost:8080/hello.HelloService/SayHello -d '{\"name\": \"Bob\"}'")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
