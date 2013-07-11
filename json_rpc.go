package rpcutil

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

var jsonRpcConnected = "200 Connected to JSON RPC"

type JsonRpcServer struct {
	*rpc.Server
}

func NewJsonRpcServer() *JsonRpcServer {
	return &JsonRpcServer{rpc.NewServer()}
}

func (server *JsonRpcServer) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Fatal("rpc.Serve: accept:", err.Error()) // TODO(r): exit?
		}
		go server.ServeConn(conn)
	}
}

func (server *JsonRpcServer) ServeConn(conn io.ReadWriteCloser) {
	server.ServeCodec(jsonrpc.NewServerCodec(conn))
}

func (server *JsonRpcServer) HandleHTTP(rpcPath string) {
	http.Handle(rpcPath, server)
}

func (server *JsonRpcServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "CONNECT" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "405 must CONNECT\n")
		return
	}

	var conn, _, err = w.(http.Hijacker).Hijack()
	if err != nil {
		log.Print("rpc hijacking ", req.RemoteAddr, ": ", err.Error())
		return
	}

	io.WriteString(conn, "HTTP/1.0 "+jsonRpcConnected+"\n\n")

	server.ServeCodec(jsonrpc.NewServerCodec(conn))
}

func DialJsonRpc(network, address, path string) (*rpc.Client, error) {
	var err error

	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	io.WriteString(conn, "CONNECT "+path+" HTTP/1.0\n\n")

	// Require successful HTTP response
	// before switching to RPC protocol.
	resp, err := http.ReadResponse(bufio.NewReader(conn), &http.Request{Method: "CONNECT"})
	if err == nil && resp.Status == jsonRpcConnected {
		return jsonrpc.NewClient(conn), nil
	}

	if err == nil {
		err = errors.New("unexpected HTTP response: " + resp.Status)
	}

	conn.Close()

	return nil, &net.OpError{
		Op:   "dial-http",
		Net:  network + " " + address,
		Addr: nil,
		Err:  err,
	}
}
