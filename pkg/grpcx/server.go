package grpcx

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	*grpc.Server
	Addr string
}

func (s *Server) Serve() {

	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("grpc service listening....\n")

	err = s.Server.Serve(l)
	if err != nil {
		panic(err)
	}

}
