package router

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/protobuf/reflect/protoreflect"
)

func init() {
	ctx := context.Background()
	if err := loadRoutes(ctx); err != nil {
		log.Fatal(err)
	}
}

var routes = make(map[string]protoreflect.MethodDescriptor)

func GetRouteInfo(urlPath string) (protoreflect.MethodDescriptor, error) {
	info, ok := routes[urlPath]
	if !ok {
		return info, fmt.Errorf("Method not found: %s", urlPath)
	}

	return info, nil
}
