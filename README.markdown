# pianobarproxy
--
Command pianobarproxy is a simple SOCKS5 shim for pianobar.

This lets you proxy pianobar through ssh (or any SOCKS5 provider).

    # Start your SOCKS5 proxy via ssh:
    ssh -v -D localhost:9080 -C -N example.com

    # Start pianobarproxy:
    pianobarproxy -socks5 :9080

    # Add the following to $HOME/.config/pianobar/config:
    proxy = http://localhost:9090

### Install

    go get github.com/robertkrimen/pianobarproxy

(http://golang.org/doc/install)

### Usage

    Usage of pianobarproxy:
        -listen="localhost:9090": The listening address
        -socks5="localhost:1080": The address of the SOCKS5 proxy

### Piano-proxy.sh
    Additional script not part of upstream.
    Starts/stops an Amazon EC2 instance on demand, then starts pianobarproxy
    Needs an existing EC2 instance. 
    Usage is ./piano-proxy.sh {start|stop}
