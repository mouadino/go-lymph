package lymph

import (
	"time"

	"github.com/mouadino/go-nano"
	"github.com/mouadino/go-nano/client"
	"github.com/mouadino/go-nano/discovery"
	"github.com/mouadino/go-nano/handler"
	"github.com/mouadino/go-nano/serializer"
	"github.com/samuel/go-zookeeper/zk"
)

var ZookeeperAnnouncer = discovery.CustomZooKeeperAnnounceResolver(
	[]string{"127.0.0.1:2181"},
	"lymph/services",
	3*time.Second,
	zk.PermAll,
	serializer.JSONSerializer{},
)

func Server(svc interface{}) *nano.Server {
	zqTrans := NewZeroMQTransport("tcp://127.0.0.1:5555")
	return nano.CustomServer(
		handler.Reflect(svc),
		zqTrans,
		NewLymphProtocol(zqTrans, serializer.MsgPackSerializer{}),
		//middleware.NewRecoverMiddleware(log.New(), true, 8*1024),
		//middleware.NewTraceMiddleware(),
		//middleware.NewLoggerMiddleware(log.New()),
	)
	//zqTrans.Stop()
}

func Client(endpoint string) nano.Client {
	// TODO: Listen on external connection.
	zqTrans := NewZeroMQTransport("tcp://127.0.0.1:5556")
	return nano.CustomClient(
		endpoint,
		NewLymphProtocol(zqTrans, serializer.MsgPackSerializer{}),
		client.NewTimeoutExt(3*time.Second),
	)
}
