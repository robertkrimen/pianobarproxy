/*
Command pianobarproxy is a simple SOCKS5 shim for pianobar.

This lets you proxy pianobar through ssh (or any SOCKS5 provider).

    # Start your SOCKS5 proxy via ssh:
    ssh -v -D localhost:9080 -C -N example.com

    # Start pianobarproxy:
    pianobarproxy -socks5 :9080

    # Add the following to $HOME/.config/pianobar/config:
    proxy = http://localhost:9090

Install

     go get github.com/robertkrimen/pianobarproxy

(http://golang.org/doc/install)

Usage

    Usage of pianobarproxy:
        -listen="localhost:9090": The listening address
        -socks5="localhost:1080": The address of the SOCKS5 proxy

*/
package main

import (
	"golang.org/x/net/proxy"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

var (
	transport *http.Transport
	dialer    proxy.Dialer
)

var (
	flag_listen = flag.String("listen", "localhost:9090", "The listening address")
	flag_socks5 = flag.String("socks5", "localhost:1080", "The address of the SOCKS5 proxy")
)

func pipe(a io.ReadWriteCloser, b io.ReadWriteCloser) {
	io.Copy(a, b)
	a.Close()
	b.Close()
}

func copyHeader(dst, src http.Header) {
	for key, entry := range src {
		for _, value := range entry {
			dst.Add(key, value)
		}
	}
}

func httpProxy(writer http.ResponseWriter, request *http.Request) {

	proxyRequest := new(http.Request)
	*proxyRequest = *request

	log.Printf("request = %s %s", request.Method, request.URL.Host)

	if strings.ToUpper(proxyRequest.Method) == "CONNECT" {
		hostPort := request.URL.Host
		pandora, err := dialer.Dial("tcp", hostPort) // tuner.pandora.com:443
		if err != nil {
			log.Printf("pianobarproxy: error: %v", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		client, writer, err := writer.(http.Hijacker).Hijack()
		writer.WriteString("HTTP/1.0 200 Connection Established\r\n\r\n")
		writer.Flush()
		go pipe(client, pandora)
		go pipe(pandora, client)
		return
	}
	proxyRequest.Proto = "HTTP/1.1"
	proxyRequest.ProtoMajor = 1
	proxyRequest.ProtoMinor = 1
	proxyRequest.Close = false

	// Remove the connection header to the backend.  We want a
	// persistent connection, regardless of what the client sent
	// to us.
	if proxyRequest.Header.Get("Connection") != "" {
		proxyRequest.Header = make(http.Header)
		copyHeader(proxyRequest.Header, request.Header)
		proxyRequest.Header.Del("Connection")
	}

	response, err := transport.RoundTrip(proxyRequest)
	if err != nil {
		log.Printf("pianobarproxy: error: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	copyHeader(writer.Header(), response.Header)

	writer.WriteHeader(response.StatusCode)

	if response.Body != nil {
		io.Copy(io.Writer(writer), response.Body)
	}
}

func getHost(target, defaultHost, defaultPort string) (string, error) {
	host, port, err := net.SplitHostPort(target)
	if err != nil {
		return "", err
	}
	if host == "" {
		host = defaultHost
	}
	if port == "" {
		port = defaultPort
	}
	return host + ":" + port, nil
}

func main() {
	flag.Parse()

	socks5, err := getHost(*flag_socks5, "localhost", "1080")
	if err != nil {
		panic(err)
	}

	log.Printf("SOCKS5/--socks5 = %s", socks5)

	{
		var err error
		dialer, err = proxy.SOCKS5("tcp", socks5, nil, &net.Dialer{})
		if err != nil {
			panic(err)
		}
	}

	transport = &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		},
	}

	listen, err := getHost(*flag_listen, "localhost", "9090")
	if err != nil {
		panic(err)
	}

	log.Printf("listen/--listen = %s", listen)
	{
		fmt.Printf("%s\n\n    %s = http://%s\n\n",
			"# Add the following to $HOME/.config/pianobar/config:",
			"proxy", listen,
		)
	}

	log.Printf("pianobarproxy")

	log.Fatal(http.ListenAndServe(listen, http.HandlerFunc(httpProxy)))
}
