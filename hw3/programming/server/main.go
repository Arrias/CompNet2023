package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

var serverDirectory = flag.String("dir", "", "root directory for server")
var concurrencyLvl = flag.Int("cl", 3, "max number parallel connections")

func log(msg string) {
	fmt.Printf("[%s]: %s\n", time.Now(), msg)
}

func process(conn net.Conn, tickets chan struct{}) {
	defer conn.Close()
	defer func() {
		tickets <- struct{}{}
	}()

	var request *http.Request
	var err error

	if request, err = http.ReadRequest(bufio.NewReader(conn)); err != nil {
		log(fmt.Sprintf("%v", err))
		return
	}

	log(request.URL.String())
	var content []byte
	if content, err = os.ReadFile(*serverDirectory + request.URL.String()); err != nil {
		log(fmt.Sprintf("%v", err))
		return
	}

	response := http.Response{Body: io.NopCloser(bytes.NewReader(content))}
	var dumpResponse []byte
	if dumpResponse, err = httputil.DumpResponse(&response, true); err != nil {
		log(fmt.Sprintf("%v", err))
		return
	}

	conn.Write(dumpResponse)
}

func main() {
	log("Start server")
	flag.Parse()

	var ln *net.TCPListener
	var err error
	tickets := make(chan struct{}, *concurrencyLvl)

	// добавляем "билеты" для соединений
	for i := 0; i < *concurrencyLvl; i++ {
		tickets <- struct{}{}
	}

	// начинаем слушать некоторый порт
	if ln, err = net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8080}); err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("error has occurred: %v\n", err)
		}
		// убеждаемся, что есть доступный билет
		<-tickets
		go process(conn, tickets)
	}
}
