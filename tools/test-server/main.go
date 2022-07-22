package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

func main() {
	println("hello world")

	addrStr := ":8080"
	address, err := net.ResolveTCPAddr("tcp", addrStr)
	if err != nil {
		log.Fatalf("could not resolve tcp address %s: %s", addrStr, err.Error())
	}

	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		log.Fatalf("could not listen on address %s: %s", address.String(), err.Error())
	}

	var handler http.HandlerFunc
	handler = func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("method: %s, content length: %d\n", r.Method, r.ContentLength)

		// handle request body
		bs := make([]byte, 1000)
		rc := r.Body
		sum := 0
		defer rc.Close()
		for {
			n, err := rc.Read(bs)

			sum += n
			fmt.Printf("read %d bytes\n", n)

			if err != nil {
				if !errors.Is(err, io.EOF) {
					w.WriteHeader(http.StatusBadRequest)
				}
				break
			}
		}

		fmt.Printf("read %d total bytes\n", sum)

		// write empty response body
		_, _ = w.Write([]byte(""))
	}

	http.Serve(listener, handler)
}
