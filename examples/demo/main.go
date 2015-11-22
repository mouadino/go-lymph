package main

import (
	"fmt"
	"time"

	"github.com/mouadino/go-lymph"
	nano "github.com/mouadino/go-nano"
	"github.com/mouadino/go-nano/client"
	"github.com/mouadino/go-nano/serializer"
)

//var echo = nano.DefaultClient("upper")
var zqTrans = lymph.NewZeroMQTransport("tcp://127.0.0.1:5556")
var echo = nano.CustomClient(
	"tcp://127.0.0.1:5555",
	lymph.NewLymphProtocol(zqTrans, serializer.MsgPackSerializer{}),
	client.NewTimeoutExt(3*time.Second),
)

func main() {
	c := time.Tick(1 * time.Second)
	i := 0
	for _ = range c {
		text := fmt.Sprintf("foo_%d", i)
		result, err := echo.Call("upper", text)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		} else {
			fmt.Printf("%s\n", result.(string))
		}
		i++
	}
}
