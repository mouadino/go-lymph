package lymph

// TODO: Extend msgpack.
// TODO: Add remote error support
// TODO: Test with other types: decimal, datetimes, unicode ...

import (
	"fmt"
	"strings"

	"github.com/mouadino/go-nano/header"
	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/serializer"
	"github.com/mouadino/go-nano/transport"
)

type lymphResponseWriter struct {
	transRW transport.ResponseWriter
	proto   *lymphProtocol
	header  header.Header
}

func (rw *lymphResponseWriter) Write(body interface{}) error {
	fmt.Printf("Sending reply to %s\n", body)
	reply, err := rw.createMessage(body, "REP")
	if err != nil {
		return err
	}
	return rw.transRW.Write(reply)
}

func (rw *lymphResponseWriter) createMessage(body interface{}, typ string) ([]string, error) {
	b, err := rw.proto.serial.Encode(body)
	if err != nil {
		return []string{}, err
	}
	header := "" // TODO
	return []string{
		typ,
		"",
		header,
		string(b),
	}, nil
}

func (rw *lymphResponseWriter) WriteError(err error) error {
	// TODO:
	fmt.Printf("Sending errros %s\n", err)
	// TODO: How to send errors ?
	reply, err := rw.createMessage("", "NACK")
	if err != nil {
		return err
	}
	return rw.transRW.Write(reply)
}

func (rw *lymphResponseWriter) Header() header.Header {
	return rw.header
}

type lymphProtocol struct {
	serial serializer.Serializer
	trans  transport.Transport
}

func NewLymphProtocol(trans transport.Transport, serial serializer.Serializer) *lymphProtocol {
	return &lymphProtocol{
		serial: serial,
		trans:  trans,
	}
}

func (proto *lymphProtocol) ReceiveRequest() (protocol.ResponseWriter, *protocol.Request) {
	transReq := <-proto.trans.Receive()

	body := transReq.Body.([]string)

	// TODO: msg_type := string(body[2])
	method := string(body[1])
	if strings.Contains(method, ".") {
		method = strings.Split(method, ".")[1]
	}

	var header map[string]string
	err := proto.serial.Decode([]byte(body[2]), &header)
	if err != nil {
		fmt.Printf("error %s\n", err)
		// TODO: Need to return an error.
		return nil, nil
	}

	var params map[string]interface{}
	err = proto.serial.Decode([]byte(body[3]), &params)
	if err != nil {
		fmt.Printf("error %s\n", err)
		// TODO: Need to return an error.
		return nil, nil
	}

	req := &protocol.Request{
		Method: method,
		Header: header,
		Params: params,
	}

	rw := &lymphResponseWriter{
		transRW: transReq.Resp,
		proto:   proto,
		header:  map[string]string{},
	}

	fmt.Printf("%v\n", req)

	return rw, req
}

func (proto *lymphProtocol) SendRequest(endpoint string, req *protocol.Request) (interface{}, error) {
	params, err := proto.serial.Encode(req.Params)
	if err != nil {
		return nil, err
	}
	headers, err := proto.serial.Encode(req.Header)
	if err != nil {
		return nil, err
	}

	msg := []string{"REQ", req.Method, string(headers), string(params)}
	// TODO: Choose a better delimiter \n\r ?
	body, err := proto.trans.Send(endpoint, []byte(strings.Join(msg, MultipartDelimeter)))
	if err != nil {
		return nil, err
	}

	frames := strings.Split(string(body), MultipartDelimeter)

	var reply interface{}
	err = proto.serial.Decode([]byte(frames[3]), &reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}
