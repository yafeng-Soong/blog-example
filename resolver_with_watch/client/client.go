package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"resolver-with-watch/register"
	"syscall"
)

var (
	port          = flag.Int("port", 8080, "")
	registerAddr  = flag.String("register", "localhost:2379", "")
	connCloseFunc []func()
)

func init() {
	flag.Parse()
}

func main() {
	defer releaseSource()

	register.InitRegister(*registerAddr)
	defer register.CloseRegister()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		log.Printf("serve at :%d", *port)
		http.HandleFunc("/hello", sayHello)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
			log.Println(err.Error())
			close(signalChan)
		}
	}()

	<-signalChan
}

func releaseSource() {
	for _, f := range connCloseFunc {
		f()
	}
}
