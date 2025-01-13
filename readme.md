## Simple Proxy

This is a simple web proxy written using raw TCP sockets. It is designed to use minimal memory by not buffering the data.

The `proxy` package is the core of the project.
The `tcp` package have all methods needed to create TCP file descriptors

## Usage

First, build the project:
```bash
go build
```

Then, run the proxy:

```bash
./simple-proxy -e :<port>
```

And that's it!
Now you can configure your browser to use the proxy, or use curl:

```bash
curl -x localhost:<port> http://www.google.com
```
