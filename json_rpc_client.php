<?php
class JsonRpcClient
{
	private $host;
	private $port;
	private $path;
	private $conn;
	private $reqId;

	function __construct($host, $port, $path) {
		$this->host = $host;
		$this->port = $port;
		$this->path = $path;
		$this->conn = NULL;
		$this->reqId = 1;
	}

	function Dial() {
		if ($this->host == "127.0.0.1" || $this->host == "localhost") 
			$host = sprintf("%s:%u", $this->host, $this->port);
		else
			$host = $this->host;

		$conn = @fsockopen($host, $this->port, $errno, $errstr, 5);

		if (!$conn) {
			return "$errstr ($errno)";
		} else {
			@fwrite($conn, "CONNECT ".$this->path." HTTP/1.0\n\n");

			stream_set_timeout($conn, 0, 3000);

			$line = @fgets($conn);

			if ($line != "HTTP/1.0 200 Connected to JSON RPC\n") {
				@fclose($conn);
				return "unexpected HTTP response: $line";
			}

			$this->conn = $conn;
		}

		return NULL;
	}

	function Call($method, $params) {
		if ($this->conn == NULL)
			return "Plaeas call Dial() first";

		$request = array(
			'method' => $method,
			'params' => array($params),
			'id' => $this->reqId,
		);

		$request = json_encode($request);

		$err = fwrite($this->conn, $request."\n");

		if ($err === false)
			return "send data failed";

		for (;;) {
			$line = @fgets($this->conn);

			if ($line != "\n") {
				break;
			}
		}

		$this->reqId += 1;

		return json_decode($line);
	}
}
?>