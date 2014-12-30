package jsonrpc

import (
	"github.com/funny/unitest"
	"log"
	"net"
	"net/http"
	"os/exec"
	"testing"
)

type Arith int

type Args struct {
	A, B int
}

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func init() {
	var server = NewServer()

	server.Register(new(Arith))
	server.HandleHTTP("/test/")

	go func() {
		var err = http.ListenAndServe("0.0.0.0:12345", nil)
		if err != nil {
			log.Fatal("Serve HTTP:", err)
		}
	}()

	go func() {
		lis, err := net.Listen("tcp", "0.0.0.0:23456")
		if err != nil {
			log.Fatal("Server TCP:", err)
		}
		server.Accept(lis)
	}()
}

func Test_Go_HTTP(t *testing.T) {
	var client, err = DialHTTP("tcp", "127.0.0.1:12345", "/test/")
	if err != nil {
		log.Fatal("Dialing:", err)
	}

	var reply int
	err = client.Call("Arith.Multiply", &Args{7, 8}, &reply)

	unitest.NotError(t, err)
	unitest.Pass(t, reply == 56)
}

func Test_Go_TCP(t *testing.T) {
	var client, err = Dial("tcp", "127.0.0.1:23456")
	if err != nil {
		log.Fatal("Dialing:", err)
	}

	var reply int
	err = client.Call("Arith.Multiply", &Args{7, 8}, &reply)

	unitest.NotError(t, err)
	unitest.Pass(t, reply == 56)
}

func Test_PHP_HTTP(t *testing.T) {
	php, err := exec.LookPath("php")
	if err != nil {
		t.Log("PHP command not found")
		return
	}

	result, err := exec.Command(php, "jsonrpc_test.php").Output()

	unitest.NotError(t, err)
	unitest.Pass(t, string(result) == "56")
}
