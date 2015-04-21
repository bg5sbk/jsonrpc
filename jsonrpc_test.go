package jsonrpc

import (
	"errors"
	"github.com/funny/unitest"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strings"
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

func (t *Arith) GetError(args *Args, reply *int) error {
	return errors.New("Error!")
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

func PHP(code string) (string, error) {
	php, err := exec.LookPath("php")
	if err != nil {
		return "", errors.New("PHP command not found")
	}
	cmd := exec.Command(php)
	cmd.Stdin = strings.NewReader("<?php " + code + ` ?>`)
	result, err := cmd.Output()
	return string(result), err
}

func Test_PHP_HTTP(t *testing.T) {
	result, err := PHP(`
		include 'jsonrpc.php';
		$client = new JsonRPC("127.0.0.1", 12345, "/test/");
		$r = $client->Call("Arith.Multiply", array('A'=>7, 'B'=>8));
		echo $r->result;
	`)
	unitest.NotError(t, err)
	unitest.Pass(t, result == "56")
}

func Test_PHP_Error(t *testing.T) {
	rpcErr, err := PHP(`
		include 'jsonrpc.php';
		$client = new JsonRPC("127.0.0.1", 12345, "/test/");
		$r = $client->Call("Arith.GetError");
		echo $r->error;
	`)
	unitest.NotError(t, err)
	unitest.Pass(t, rpcErr == "Error!")
}
