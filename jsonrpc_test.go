package jsonrpc

import (
	"github.com/funny/unitest"
	"log"
	"net/http"
	"testing"
)

type TestAPI int

type Args struct {
	A, B int
}

func (t *TestAPI) Hello(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func init() {
	var server = NewServer()

	server.Register(new(TestAPI))
	server.HandleHTTP("/test/")

	go func() {
		var err = http.ListenAndServe("0.0.0.0:12345", nil)
		if err != nil {
			log.Fatal("Serve Http:", err)
		}
	}()
}

func Test_JsonRPC(t *testing.T) {
	var client, err = DialHTTP("tcp", "127.0.0.1:12345", "/test/")
	if err != nil {
		log.Fatal("Dialing:", err)
	}

	var reply int
	err = client.Call("TestAPI.Hello", &Args{7, 8}, &reply)

	unitest.NotError(t, err)
	unitest.Pass(t, reply == 56)
}
