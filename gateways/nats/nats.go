package nats

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
)

var Module = fx.Provide(New)

type Gateway interface {
	GetConnection() *nats.Conn
}

type gateway struct {
	Connection *nats.Conn
}

func New() Gateway {
	nc, err := nats.Connect("nats://nats-server:4222")
	if err != nil {
		panic(err)
	}

	return &gateway{
		Connection: nc,
	}
}

func (g *gateway) GetConnection() *nats.Conn {
	return g.Connection
}
