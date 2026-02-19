package rpc

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"tiny-gateway/resolver"
	"tiny-gateway/router"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/dynamicpb"
)

// func init() {
// 	var err error
// 	conn, err = grpc.NewClient(
// 		"192.168.0.122:46509",
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 	)
// 	if err != nil {
// 		log.Fatalf("create gRPC client: %v", err)
// 	}
// }

var connMap sync.Map

func getClient(serviceName string) (*grpc.ClientConn, error) {
	val, ok := connMap.Load(serviceName)
	if ok {
		return val.(*grpc.ClientConn), nil
	}

	target := fmt.Sprintf("%s://%s", resolver.Scheme, serviceName)
	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		return nil, err
	}

	if val, ok := connMap.Load(serviceName); ok {
		return val.(*grpc.ClientConn), nil
	}

	connMap.Store(serviceName, conn)
	return conn, nil
}

func HandleRPC(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only supports POST now", http.StatusMethodNotAllowed)
		return
	}

	// find the matching gRPC method
	// urlPath is like `/echo/echo.EchoService/Echo`
	methodDesc, err := router.GetRouteInfo(r.URL.Path)
	if err != nil {
		http.Error(w, fmt.Sprintf("Method not found: %s", r.URL.Path), http.StatusNotFound)
		return
	}

	urlPath := strings.SplitN(r.URL.Path, "/", 3)
	if len(urlPath) < 3 {
		http.Error(w, "Invalide URL path", http.StatusBadRequest)
		return
	}

	serviceName, methodPath := urlPath[1], urlPath[2]
	conn, err := getClient(serviceName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect server: %v", err), http.StatusBadGateway)
		return
	}

	// read the HTTP JSON body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// JSON -> Protobuf
	reqMsg := dynamicpb.NewMessage(methodDesc.Input())
	if err := protojson.Unmarshal(body, reqMsg); err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal request: %v", err), http.StatusBadRequest)
		return
	}

	// call the gRPC method
	respMsg := dynamicpb.NewMessage(methodDesc.Output())
	if err := conn.Invoke(r.Context(), methodPath, reqMsg, respMsg); err != nil {
		http.Error(w, fmt.Sprintf("Failed to call gRPC: %v", err), http.StatusBadGateway)
		return
	}

	// Protobuf -> JSON bytes
	jsonBytes, err := protojson.Marshal(respMsg)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal response: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}
