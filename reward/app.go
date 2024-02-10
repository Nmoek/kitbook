package main

import (
	"kitbook/internal/events"
	"kitbook/pkg/grpcx"
)

type App struct {
	consumers []events.Consumer
	rpcServer *grpcx.Server
}
