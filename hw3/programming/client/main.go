package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var host = flag.String("host", "127.0.0.1", "host")
var port = flag.String("port", "8080", "port")
var filename = flag.String("fn", "main.go", "filename")

func main() {
	flag.Parse()
	var err error
	var conn net.Conn
	if conn, err = net.Dial("tcp", *host+":"+*port); err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	request := http.Request{Method: "GET", URL: &url.URL{Host: *host + ":" + *port, Path: "/" + *filename}}
	var dummyRequest []byte
	if dummyRequest, err = httputil.DumpRequest(&request, false); err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	conn.Write(dummyRequest)
	var response *http.Response
	if response, err = http.ReadResponse(bufio.NewReader(conn), &request); err != nil {
		fmt.Printf("%v", err)
		return
	}

	bodyReader := bufio.NewReader(response.Body)
	for {
		line, err := bodyReader.ReadString('\n')
		if errors.Is(err, io.EOF) {
			break
		}
		fmt.Println(line)
	}
}
