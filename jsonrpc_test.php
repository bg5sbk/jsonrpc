<?php
include 'jsonrpc.php';

$client = new JsonRpcClient("127.0.0.1", 12345, "/test/");

echo $client->Dial();
echo "\n";

var_export($client->Call("Arith.Multiply", array('A'=>7, 'B'=>8)));
echo "\n";

var_export($client->Call("Arith.Multiply", array('A'=>6, 'B'=>6)));
echo "\n";
?>