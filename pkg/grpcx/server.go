package grpcx

import (
	"context"
	"fmt"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"google.golang.org/grpc"
	"kitbook/pkg/logger"
	"kitbook/pkg/netx"
	"net"
	"strconv"
	"time"
)

type Server struct {
	*grpc.Server
	cli      *etcdv3.Client //保持依赖注入 或 独立创建
	EtcdAddr string
	Port     int
	kaCancel func()
	BizName  string
	L        logger.Logger
}

func (s *Server) Serve() {
	addr := "localhost" + ":" + strconv.Itoa(s.Port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	err = s.register()
	if err != nil {
		s.Close()
		panic(err)
	}

	fmt.Printf("grpc service listening....\n")
	err = s.Server.Serve(l)
	if err != nil {
		panic(err)
	}

}

func (s *Server) Close() error {
	if s.kaCancel != nil {
		// 关闭续约、心跳
		s.kaCancel()
	}

	if s.cli != nil {
		// 关闭服务注册连接
		return s.cli.Close()
	}

	// 关闭gRPC服务
	s.GracefulStop()

	return nil
}

func (s *Server) register() error {

	// 创建服务注册客户端
	cli, err := etcdv3.NewFromURL(s.EtcdAddr)
	if err != nil {
		return err
	}
	s.cli = cli

	em, err := endpoints.NewManager(s.cli, "service/"+s.BizName)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 初始租约
	ttl := int64(5)
	leaseResp, err := s.cli.Grant(ctx, ttl)
	if err != nil {
		return err
	}

	// 向etcd注册服务
	addr := netx.GetOutboundIP() + ":" + strconv.Itoa(s.Port) //服务对外IP
	key := "service/" + s.BizName + "/" + addr                // 组装服务唯一的key
	err = em.AddEndpoint(ctx, key, endpoints.Endpoint{
		Addr: addr,
	}, etcdv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}

	// 开启续约
	kaCtx, kaCancel := context.WithCancel(context.Background())
	s.kaCancel = kaCancel
	ch, err := s.cli.KeepAlive(kaCtx, leaseResp.ID)
	go func() {
		for kaResp := range ch {
			s.L.INFO("interactive服务续约",
				logger.Field{"msg", kaResp.String()})
		}
	}()

	return err

}
