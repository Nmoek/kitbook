package ioc

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	followv1 "kitbook/api/proto/gen/follow/v1"
)

func InitFollowClient(cli *etcdv3.Client) followv1.FollowServiceClient {

	type Config struct {
		Addr   string `yaml:"addr"`
		Secure bool
	}

	var cfg Config

	err := viper.UnmarshalKey("grpc.client.follow", &cfg)
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
		cfg = Config{}
		err2 := viper.UnmarshalKey("grpc.client.follow", &cfg)
		if err2 != nil {
			panic(err2)
		}
	})

	return followv1.NewFollowServiceClient(cc)
}
