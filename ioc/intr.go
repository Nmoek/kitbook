package ioc

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	intrv1 "kitbook/api/proto/gen/intr/v1"
	"kitbook/interactive/service"
	"kitbook/internal/client"
)

// @func: InitIntrClientV1
// @date: 2024-01-17 01:21:15
// @brief: 注册中心-只发起远程调用
// @author: Kewin Li
// @param svc
// @return intrv1.InteractiveServiceClient
func InitIntrClientV1(svc service.InteractiveService, cli *etcdv3.Client) intrv1.InteractiveServiceClient {
	type Config struct {
		Addr      string `yaml:"addr"`
		Secure    bool
		Threshold int32
	}

	var cfg Config
	err := viper.UnmarshalKey("grpc.client.intr", &cfg)
	if err != nil {
		panic(err)
	}

	// 创建远程服务发现客户端
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
		err2 := viper.UnmarshalKey("grpc.client.intr", &cfg)
		if err2 != nil {
			panic(err2)
		}

	})

	return intrv1.NewInteractiveServiceClient(cc)

}

// @func: InitIntrClient
// @date: 2024-01-17 01:20:53
// @brief: 做灰度控制, 本地调用+远程调用
// @author: Kewin Li
// @param svc
// @return intrv1.InteractiveServiceClient
func InitIntrClient(svc service.InteractiveService) intrv1.InteractiveServiceClient {
	type Config struct {
		Addr      string `yaml:"addr"`
		Secure    bool
		Threshold int32
	}

	var cfg Config
	err := viper.UnmarshalKey("grpc.client.intr", &cfg)
	if err != nil {
		panic(err)
	}

	var opts []grpc.DialOption
	if !cfg.Secure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	cc, err := grpc.Dial(cfg.Addr, opts...)
	if err != nil {
		panic(err)
	}

	remote := intrv1.NewInteractiveServiceClient(cc)
	local := client.NewLocalInteractiveServiceAdapter(svc)
	res := client.NewInteractiveClient(remote, local)

	viper.OnConfigChange(func(in fsnotify.Event) {
		cfg = Config{}
		err2 := viper.UnmarshalKey("grpc.client.intr", &cfg)
		if err2 != nil {
			panic(err2)
		}

		// 流量分流阈值
		res.UpdateThreshold(cfg.Threshold)
	})

	return res

}
