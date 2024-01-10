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

<<<<<<< HEAD
	fmt.Printf("grpc service listening....\n")
=======
	fmt.Printf("grpc service listening....")
>>>>>>> 36b0fdf641cfc6910a7686ac6697827f8c0a2482

	err = s.Server.Serve(l)
	if err != nil {
		panic(err)
	}

}
