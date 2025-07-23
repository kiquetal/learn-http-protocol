# Learn HTTP Protocol

This project is a collection of simple Go programs to learn about the HTTP protocol.

## Index

* [Commands](#commands)
  * [HTTP Server](#http-server)
  * [TCP Listener](#tcp-listener)
  * [UDP Sender](#udp-sender)
* [Interaction Diagram](#interaction-diagram)


## Commands

### HTTP Server

This command starts a simple HTTP server on port 42069.

**Usage:**

```bash
go run ./cmd/httpserver
```

The server has the following endpoints:

*   `/`: Returns a 200 OK response with a simple HTML page.
*   `/video`: Returns a video file.
*   `/httpbin/*`: Forwards the request to `https://httpbin.org` and returns the response.
*   `/yourproblem`: Returns a 400 Bad Request response.
*   `/myproblem`: Returns a 500 Internal Server Error response.

### TCP Listener

This command starts a TCP listener on port 42069 and prints the incoming requests.

**Usage:**

```bash
go run ./cmd/tcplistener
```

### UDP Sender

This command sends UDP packets to port 42069.

**Usage:**

```bash
go run ./cmd/udpsender
```

After running the command, you can type any message and press Enter to send it to the TCP listener.

## Interaction Diagram

```ascii
+-----------------+      +-----------------+      +-----------------+
|   UDP Sender    |----->|  TCP Listener   |<-----|   HTTP Client   |
+-----------------+      +-----------------+      +-----------------+
                           ^
                           |
                           v
+-----------------+      +-----------------+
|  HTTP Server    |----->|   httpbin.org   |
+-----------------+      +-----------------+
```
