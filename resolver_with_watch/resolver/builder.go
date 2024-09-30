package resolver

import (
	"context"
	"fmt"
	"log"
	"resolver-with-watch/register"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
)

const (
	Scheme    = "costum"
	logPrefix = "[resolver] "
)

var (
	cancel      context.CancelFunc
	serviceInfo *serviceInfoCache
	wacthOver   <-chan bool
)

func init() {
	serviceInfo = &serviceInfoCache{cache: make(map[string][]resolver.Address)}
	builder := manual.NewBuilderWithScheme(Scheme)
	builder.BuildCallback = buildCallbackfunc
	builder.CloseCallback = closeCallBack
	resolver.Register(builder)
}

// store addresses information of all kind of services, not noly hello-server.
type serviceInfoCache struct {
	cache map[string][]resolver.Address
}

func (s *serviceInfoCache) GetAddresses(serviceName string) []resolver.Address {
	addrs, ok := s.cache[serviceName]
	if !ok {
		addrs = []resolver.Address{}
	}

	return addrs
}

func (s *serviceInfoCache) UpdateAddress(serviceName string, addrs []resolver.Address) {
	s.cache[serviceName] = addrs
}

func sliceToString(addrs []resolver.Address) string {
	var tmp []string
	for _, addr := range addrs {
		tmp = append(tmp, addr.Addr)
	}
	return fmt.Sprint(tmp)
}

func buildCallbackfunc(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) {
	serviceName := target.URL.Host
	hosts := register.QueryAddress(serviceName)
	log.Printf("%sfound service %s in %v ", logPrefix, serviceName, sliceToString(hosts))

	serviceInfo.UpdateAddress(serviceName, hosts)
	addrs := serviceInfo.GetAddresses(serviceName)
	if err := cc.UpdateState(resolver.State{Addresses: addrs}); err != nil {
		log.Printf("%supdateState error: %s\n", logPrefix, err.Error())
	}

	var ctx context.Context
	ctx, cancel = context.WithCancel(context.Background())
	wacthOver = register.WatchAddress(ctx, cc, serviceName, watchCallBack)
}

func closeCallBack() {
	if cancel != nil {
		cancel()
		<-wacthOver
		log.Println(logPrefix, "closed")
	}
}

func watchCallBack(cc resolver.ClientConn, serviceName string, ev *clientv3.Event) {
	addrs := serviceInfo.GetAddresses(serviceName)
	newAddrs := updateAddresses(addrs, ev)
	serviceInfo.UpdateAddress(serviceName, newAddrs)
	log.Printf("%sservice %s update addrs: %s", logPrefix, serviceName, sliceToString(newAddrs))
	err := cc.UpdateState(resolver.State{Addresses: newAddrs})
	if err != nil {
		log.Printf("%supdateState error: %s", logPrefix, err.Error())
	}
}

func updateAddresses(addrs []resolver.Address, ev *clientv3.Event) []resolver.Address {
	addrMap := make(map[resolver.Address]struct{})
	for _, addr := range addrs {
		addrMap[addr] = struct{}{}
	}

	switch ev.Type {
	case mvccpb.PUT:
		key := resolver.Address{Addr: string(ev.Kv.Value)}
		addrMap[key] = struct{}{}
	case mvccpb.DELETE:
		key := resolver.Address{Addr: string(ev.PrevKv.Value)}
		delete(addrMap, key)
	}

	var res []resolver.Address
	for addr := range addrMap {
		res = append(res, addr)
	}

	return res
}
