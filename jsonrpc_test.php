<?php
include 'jsonrpc.php';

$client = new JsonRPC("127.0.0.1", 12345, "/test/");
$r = $client->Call("Arith.Multiply", array('A'=>7, 'B'=>8));
echo $r->result;
?>