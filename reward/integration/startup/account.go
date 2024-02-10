package startup

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	accountv1 "kitbook/api/proto/gen/account/v1"
)

const accountSvcKey = "grpc.client.account"

func InitAccountClient(cli *etcdv3.Client) accountv1.AccountServiceClient {
	type Config struct {
		Addr   string `yaml:"addr"`
		Secure bool
	}
	var cfg Config
	err := viper.UnmarshalKey(accountSvcKey, &cfg)
	if err != nil {
		panic(err)
	}

	resolverEtcd, err := resolver.NewBuilder(cli)
	if err != nil {
		panic(err)
	}

	opts := []grpc.DialOption{grpc.WithResolvers(resolverEtcd)}
	if !cfg.Secure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	cc, err := grpc.Dial(cfg.Addr, opts...)
	if err != nil {
		panic(err)
	}

	viper.OnConfigChange(func(in fsnotify.Event) {
		cfg := Config{}
		err2 := viper.UnmarshalKey(accountSvcKey, &cfg)
		if err2 != nil {
			panic(err2)
		}
	})

	return accountv1.NewAccountServiceClient(cc)
}
