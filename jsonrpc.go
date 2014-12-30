package jsonrpc

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

type Server struct {
	*rpc.Server
}

func NewServer() *Server {
	return &Server{rpc.NewServer()}
}

func (server *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Fatal("rpc.Serve: accept:", err.Error()) // TODO(r): exit?
		}
		go server.ServeConn(conn)
	}
}

func (server *Server) ServeConn(conn io.ReadWriteCloser) {
	server.ServeCodec(jsonrpc.NewServerCodec(conn))
}

func (server *Server) HandleHTTP(rpcPath string) {
	http.Handle(rpcPath, server)
}

func (server *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var conn, _, err = w.(http.Hijacker).Hijack()
	if err != nil {
		log.Print("rpc hijacking ", req.RemoteAddr, ": ", err.Error())
		return
	}
	io.WriteString(conn, "HTTP/1.1 "+jsonRpcConnected+"\n")
	//io.WriteString(conn, "Access-Control-Allow-Origin: *\n")
	io.WriteString(conn, "Content-Type: application/json\n\n")
	server.ServeCodec(jsonrpc.NewServerCodec(conn))
}

func Dial(network, address string) (*rpc.Client, error) {
	return jsonrpc.Dial(network, address)
}

func DialHTTP(network, address, path string) (*rpc.Client, error) {
	var err error

	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	io.WriteString(conn, "GET "+path+" HTTP/1.1\n\n")

	// Require successful HTTP response
	// before switching to RPC protocol.
	resp, err := http.ReadResponse(bufio.NewReader(conn), nil)
	if err == nil {
		if resp.Status == jsonRpcConnected {
			return jsonrpc.NewClient(conn), nil
		}
	}

	conn.Close()

	return nil, &net.OpError{
		Op:  "JsonRPC dial to",
		Net: network + "://" + address,
		Err: errors.New("unexpected HTTP response: " + resp.Status),
	}
}
