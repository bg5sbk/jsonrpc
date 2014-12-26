项目介绍
=======

Go语言作为面向服务端编程的语言，其标准库中自带了功能完整的RPC模块，但是这个模块用的是gob格式传输数据，当开发Go程序间通讯时是很好用的，但是要让Go之外的别的编程语言解析gob格式就很麻烦了。

夸语言通讯首选的当然是XML或者JSON之类的文本协议，而Go内置的JSON RPC模块，只实现了协议解析部分，没有实现传输的部分。要用之前还得自己实现一整套接口注册、连接处理、请求解析，实在是麻烦。

真有趣团队的项目中用到了JSON RPC，我们用了一个比较取巧的办法，尽量利用Go自带模块的代码来实现一套完整的JSON RPC服务端和客户端，并对应实现了一个PHP客户端.

效果演示
=======

下面是这个JSON RPC模块的用法示例。

GO代码如下：
```go
package main

import (
	"errors"
	"fmt"
	"github.com/funny/jsonrpc"
	"log"
	"net/http"
	"time"
	"sync"
)

type Arith int

type Args struct {
	A, B int
}

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func main() {
	var server = jsonrpc.NewServer()

	server.Register(new(Arith))
	server.HandleHTTP("/test/")

	var wg sync.WatiGroup
	wg.Add(1)
	go func() {
		wg.Done()
		var err = http.ListenAndServe("0.0.0.0:12345", nil)
		if err != nil {
			log.Fatal("Serve Http:", err)
		}
	}()

	var client, err = jsonrpc.DialHTTP("tcp", "127.0.0.1:12345", "/test/")
	if err != nil {
		log.Fatal("Dialing:", err)
	}

	var args = &Args{7, 8}
	var reply int

	err = client.Call("Arith.Multiply", args, &reply)

	if err != nil {
		log.Fatal("Call:", err)
	}

	log.Printf("Result: %d * %d = %d", args.A, args.B, reply)

	fmt.Scanln()
}
```

PHP调用端代码如下：
```php
include 'jsonrpc.php';

$client = new JsonRpcClient("127.0.0.1", 12345, "/test/");

echo $client->Dial();
echo "\n";

var_export($client->Call("Arith.Multiply", array('A'=>7, 'B'=>8)));
echo "\n";

var_export($client->Call("Arith.Multiply", array('A'=>6, 'B'=>6)));
echo "\n";
```

PHP端输出结果：

```
stdClass::__set_state(array(
   'id' => 1,
   'result' => 56,
   'error' => NULL,
))
stdClass::__set_state(array(
   'id' => 2,
   'result' => 36,
   'error' => NULL,
))
```