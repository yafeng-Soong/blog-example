package router

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/bufbuild/protocompile"
)

func loadRoutes(ctx context.Context) error {
	// 约定：proto 目录保存了其他服务的proto文件
	// proto/
	//   echo.proto
	//   foo.proto

	compiler := protocompile.Compiler{
		Resolver: &protocompile.SourceResolver{
			ImportPaths: []string{"proto"},
		},
	}

	entries, err := os.ReadDir("proto")
	if err != nil {
		return fmt.Errorf("read proto dir: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// treat dir name as app name
		appName := entry.Name()
		if err := loadAppRoutes(ctx, &compiler, appName); err != nil {
			return err
		}
	}

	return nil
}

// 加载单个 app 目录下的 proto 与路由配置
func loadAppRoutes(ctx context.Context, compiler *protocompile.Compiler, appName string) error {
	// 检查该目录下是否包含约定命名的 proto 文件
	protoFilePath := path.Join(appName, appName+".proto")
	compiled, err := compiler.Compile(ctx, protoFilePath)
	if err != nil {
		return fmt.Errorf("compile proto %s: %w", protoFilePath, err)
	}

	fileDesc := compiled.FindFileByPath(protoFilePath)
	if fileDesc == nil {
		return fmt.Errorf("file not found: %s", protoFilePath)
	}

	services := fileDesc.Services()
	for i := range services.Len() {
		svc := services.Get(i)
		methods := svc.Methods()
		for j := range methods.Len() {
			method := methods.Get(j)
			fullMethod := fmt.Sprintf("/%s/%s/%s", appName, svc.FullName(), method.Name())
			routes[fullMethod] = method
			log.Printf("registering route %s\n", fullMethod)
		}
	}

	return nil
}
