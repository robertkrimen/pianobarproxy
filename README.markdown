# pianobarproxy
--
Command pianobarproxy is a simple SOCKS5 shim for pianobar.

This lets you proxy pianobar through ssh (or any SOCKS5 provider).

    # Start your SOCKS5 proxy via ssh:
    ssh -v -D 9080 -C -N

    # Start pianobarproxy:
    pianobarproxy --remote :9080

    # Add the following to $HOME/.config/pianobar/config:
    proxy = http://localhost:9090


### Install

http://golang.org/doc/install

    go get github.com/robertkrimen/pianobarproxy

--
**godocdown** http://github.com/robertkrimen/godocdown
