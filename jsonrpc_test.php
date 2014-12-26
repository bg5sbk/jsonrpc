<?php
include 'jsonrpc.php';

$client = new JsonRpcClient("127.0.0.1", 12345, "/test/");
$client->Dial();

$r = $client->Call("Arith.Multiply", array('A'=>7, 'B'=>8));
echo $r->result;
?>