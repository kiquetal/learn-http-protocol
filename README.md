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

**Testing with `telnet`:**

You can use `telnet` to connect to the TCP listener and send messages.

1.  In one terminal, start the TCP listener:

    ```bash
    go run ./cmd/tcplistener
    ```

2.  In another terminal, connect to the listener using `telnet`:

    ```bash
    telnet localhost 42069
    ```

3.  After connecting, you can type any message and press Enter. The message will be sent to the TCP listener, and you will see the parsed request details in the listener's terminal.

### UDP Sender

This command sends UDP packets to port 42069.

**Usage:**

```bash
go run ./cmd/udpsender
```

After running the command, you can type any message and press Enter to send it as a UDP packet to port 42069.

**Testing with `netcat`:**

You can use `netcat` (`nc`) to listen for the UDP packets sent by the `udpsender`.

1.  In one terminal, start `netcat` to listen for UDP packets on port 42069:

    ```bash
    nc -u -l -p 42069
    ```

2.  In another terminal, run the UDP sender:

    ```bash
    go run ./cmd/udpsender
    ```

3.  Type a message in the `udpsender` terminal (e.g., "hello world") and press Enter. The message will appear in the `netcat` terminal.

## Interaction Diagram

The following diagram shows the interaction between the different components. The `UDP Sender` is a standalone component to demonstrate UDP communication and can be tested with a UDP listener like `netcat`.

```ascii
+-----------------+      +-----------------+      +-----------------+
|   UDP Sender    |      |  TCP Listener   |<-----|   HTTP Client   |
+-----------------+      +-----------------+      +-----------------+
                           ^
                           |
                           v
+-----------------+      +-----------------+
|  HTTP Server    |----->|   httpbin.org   |
+-----------------+      +-----------------+
```
