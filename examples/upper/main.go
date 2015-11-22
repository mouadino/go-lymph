package main

import (
	"strings"

	"github.com/mouadino/go-lymph"
	"github.com/mouadino/go-nano"
	"github.com/mouadino/go-nano/handler"
	"github.com/mouadino/go-nano/serializer"
)

type Upper struct{}

func (Upper) Upper(text string) string {
	return strings.ToUpper(text)
}

func main() {
	zqTrans := lymph.NewZeroMQTransport("tcp://127.0.0.1:5555")
	server := nano.CustomServer(
		handler.Reflect(Upper{}),
		zqTrans,
		lymph.NewLymphProtocol(zqTrans, serializer.MsgPackSerializer{}),
		//middleware.NewRecoverMiddleware(log.New(), true, 8*1024),
		//middleware.NewTraceMiddleware(),
		//middleware.NewLoggerMiddleware(log.New()),
	)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
	//zqTrans.Stop()
}
