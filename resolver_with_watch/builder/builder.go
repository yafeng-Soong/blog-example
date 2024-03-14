package builder

import (
	"context"
	"log"
	"resolver-with-watch/register"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
)

const Scheme = "costum"

var (
	cancel      context.CancelFunc
	serviceInfo *serviceInfoCache
	wacthOver   <-chan bool
)

func init() {
	serviceInfo = &serviceInfoCache{cache: make(map[string]map[string]resolver.Address)}
	builder := manual.NewBuilderWithScheme(Scheme)
	builder.BuildCallback = buildCallbackfunc
	builder.CloseCallback = closeCallBack
	resolver.Register(builder)
}

// store addresses information of all kind of services, not noly hello-server.
type serviceInfoCache struct {
	cache map[string]map[string]resolver.Address
}

func (s *serviceInfoCache) GetAddresses(serviceName string) []resolver.Address {
	var addrs []resolver.Address
	info, ok := s.cache[serviceName]
	if !ok {
		log.Println("[error] did not watch service: ", serviceName)
		return addrs
	}

	for _, k := range info {
		addrs = append(addrs, k)
	}
	return addrs
}

func (s *serviceInfoCache) SetAddresses(serviceName string, addrs map[string]resolver.Address) {
	s.cache[serviceName] = addrs
}

func (s *serviceInfoCache) UpdateAddress(serviceName string, ev *clientv3.Event) {
	info, ok := s.cache[serviceName]
	if !ok {
		log.Println("[error] did not watch service: ", serviceName)
		return
	}

	instance := string(ev.Kv.Key)
	switch ev.Type {
	case mvccpb.PUT:
		info[instance] = resolver.Address{Addr: string(ev.Kv.Value)}
	case mvccpb.DELETE:
		delete(info, instance)
	}

}

func buildCallbackfunc(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) {
	serviceName := target.URL.Host
	hosts := register.QueryAddress(serviceName)
	if len(hosts) == 0 {
		log.Fatalf("service %s not found", target)
	}

	serviceInfo.SetAddresses(serviceName, hosts)
	addrs := serviceInfo.GetAddresses(serviceName)
	err := cc.UpdateState(resolver.State{Addresses: addrs})
	if err != nil {
		log.Fatalf("UpdateState error: %s", err.Error())
	}

	// TODO watch address from etcd
	var ctx context.Context
	ctx, cancel = context.WithCancel(context.Background())
	wacthOver = register.WatchAddress(ctx, cc, serviceName, watchCallBack)
}

func closeCallBack() {
	// TODO cancel etcd watcher
	if cancel != nil {
		cancel()
	}
	<-wacthOver
	log.Println("client closed")
}

func watchCallBack(cc resolver.ClientConn, serviceName string, ev *clientv3.Event) {
	serviceInfo.UpdateAddress(serviceName, ev)
	addrs := serviceInfo.GetAddresses(serviceName)
	log.Printf("service: %s update, addrs: %v", serviceName, addrs)
	err := cc.UpdateState(resolver.State{Addresses: addrs})
	if err != nil {
		log.Printf("UpdateState error: %s", err.Error())
	}
}
