/*

Frames:
     --------------------------------------------------
    | Endpoint | correlation ID | protocol frames .... |
     --------------------------------------------------
*/

package lymph

import (
	"fmt"
	"strings"
	"time"

	"github.com/mouadino/go-nano/transport"
	"github.com/nu7hatch/gouuid"
	zmq "github.com/pebbe/zmq4"
)

// TODO: Generate a better delimeter.
const MultipartDelimeter = "\n"

type zeroMQResponseWriter struct {
	sock          *zmq.Socket
	endpoint      string
	correlationID string
}

func (rw *zeroMQResponseWriter) Write(b interface{}) error {
	fmt.Println(b)
	body := b.([]string)
	body = append([]string{rw.endpoint, rw.correlationID}, body...)
	fmt.Println("Reply ", body)
	t, err := rw.sock.SendMessage(body)

	fmt.Println("sent", t)
	return err
}

// TODO: Disconnecting
// TODO: Heartbeat.
type zeroMQTransport struct {
	endpoint        string
	sendSock        *zmq.Socket
	rcvSock         *zmq.Socket
	reqs            chan transport.Request
	connections     map[string]string
	pendingRequests map[string]chan []string
}

func NewZeroMQTransport(endpoint string) *zeroMQTransport {
	return &zeroMQTransport{
		endpoint:        endpoint,
		reqs:            make(chan transport.Request),
		connections:     make(map[string]string),
		pendingRequests: make(map[string]chan []string),
	}
}

func (trans *zeroMQTransport) Stop() {
	// TODO: Call me.
	fmt.Printf("Closing ....")
	for k, _ := range trans.connections {
		trans.sendSock.Disconnect(k)
	}
	trans.sendSock.Close()
	trans.rcvSock.Close()
}

func (trans *zeroMQTransport) Listen() error {
	var err error
	trans.rcvSock, err = zmq.NewSocket(zmq.ROUTER)
	if err != nil {
		return err
	}
	trans.rcvSock.SetIdentity(trans.endpoint)
	// TODO: Pick random port
	trans.rcvSock.Bind(trans.endpoint)

	trans.sendSock, err = zmq.NewSocket(zmq.ROUTER)
	if err != nil {
		return err
	}
	trans.sendSock.SetIdentity(trans.endpoint)

	fmt.Println("Listening ...")
	go trans.serve()
	return nil
}

func (trans *zeroMQTransport) serve() {
	for {
		msg, err := trans.rcvSock.RecvMessage(0)
		if err != nil {
			fmt.Printf("ZEROMQ error: %s", err)
			continue
		}

		endpoint := msg[0]
		msgID := msg[1]

		msg = msg[2:]
		if _, ok := trans.pendingRequests[msgID]; ok {
			fmt.Println("Reply received")
			trans.pendingRequests[msgID] <- msg
		} else {
			// TODO: Remove when disconnected.
			if _, ok := trans.connections[endpoint]; !ok {
				fmt.Println("connect to ", endpoint)
				err := trans.sendSock.Connect(endpoint)
				if err != nil {
					fmt.Printf("connection error %s", err)
				}
				time.Sleep(2 * time.Millisecond) // FIXME: Seriously :(
				trans.connections[endpoint] = ""
			}
			fmt.Println("New request ", msg)
			trans.reqs <- transport.Request{
				Body: msg,
				Resp: &zeroMQResponseWriter{
					sock:          trans.sendSock,
					endpoint:      endpoint,
					correlationID: msgID,
				},
			}
		}
	}
}

func (trans *zeroMQTransport) Addr() string {
	return trans.endpoint
}

func (trans *zeroMQTransport) Send(endpoint string, b []byte) ([]byte, error) {
	if trans.sendSock == nil {
		trans.Listen()
	}
	body := strings.Split(string(b), MultipartDelimeter)
	fmt.Println(endpoint, body)

	err := trans.sendSock.Connect(endpoint)
	if err != nil {
		return []byte{}, err
	}
	time.Sleep(2 * time.Millisecond) // FIXME: Seriously :(
	msgID, err := uuid.NewV4()
	if err != nil {
		return []byte{}, err
	}
	body = append([]string{endpoint, msgID.String()}, body...)
	_, err = trans.sendSock.SendMessage(body)
	if err != nil {
		return []byte{}, err
	}
	fmt.Println("Waiting for response")
	trans.pendingRequests[msgID.String()] = make(chan []string)
	reply := <-trans.pendingRequests[msgID.String()]
	// TODO: Remove from pendingRequests.
	return []byte(strings.Join(reply, MultipartDelimeter)), nil
}

func (trans *zeroMQTransport) Receive() <-chan transport.Request {
	return trans.reqs
}
