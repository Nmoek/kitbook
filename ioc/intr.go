package ioc

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	intrv1 "kitbook/api/proto/gen/intr/v1"
	"kitbook/interactive/service"
	"kitbook/internal/client"
)

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

	fmt.Printf("cfg= %v \n", cfg)

	viper.OnConfigChange(func(in fsnotify.Event) {
		cfg = Config{}
		err2 := viper.UnmarshalKey("grpc.client.intr", &cfg)
		if err2 != nil {
			panic(err2)
		}

		fmt.Printf("流量阈值=%v \n", cfg.Threshold)
		// 流量分流阈值
		res.UpdateThreshold(cfg.Threshold)
	})

	return res

}
