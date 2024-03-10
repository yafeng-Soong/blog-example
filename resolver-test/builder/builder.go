package builder

import (
	"log"
	"resolver-test/register"

	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
)

const Scheme = "costum"

func init() {
	builder := manual.NewBuilderWithScheme(Scheme)
	builder.BuildCallback = buildCallbackfunc
	resolver.Register(builder)
}

func buildCallbackfunc(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) {
	var address []resolver.Address
	serviceName := target.URL.Host
	hosts := register.QueryAddress(serviceName)
	if len(hosts) == 0 {
		log.Fatalf("service %s not found", target)
	}

	for _, host := range hosts {
		address = append(address, resolver.Address{
			Addr: host,
		})
	}

	err := cc.UpdateState(resolver.State{Addresses: address})
	if err != nil {
		log.Fatalf("UpdateState error: %s", err.Error())
	}
}
