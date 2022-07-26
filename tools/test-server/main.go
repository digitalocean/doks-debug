package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
)

func main() {
	addrStr := ":8080"
	address, err := net.ResolveTCPAddr("tcp", addrStr)

	logger := log.New(os.Stdout, "test-server", log.LstdFlags|log.Lmsgprefix|log.LUTC|log.Lmicroseconds)

	if err != nil {
		logger.Fatalf("could not resolve tcp address %s: %s", addrStr, err.Error())
	}

	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		logger.Fatalf("could not listen on address %s: %s", address.String(), err.Error())
	}

	var (
		handler         http.HandlerFunc
		connectionCount uint64
	)
	handler = func(w http.ResponseWriter, r *http.Request) {
		var (
			thisConnectionNum = atomic.AddUint64(&connectionCount, 1)
			prefixParts       = []string{fmt.Sprintf("%d", thisConnectionNum)}

			requestUuid = r.Header.Get("x-request-id")

			requestDigest  = r.Header.Get("digest")
			validateDigest bool

			buffer         = make([]byte, 1024+1024)
			totalBytes     int
			hasher         = sha256.New()
			responseStatus = http.StatusOK
		)

		if requestUuid != "" {
			prefixParts = append(prefixParts, requestUuid)
		}
		logger := log.New(os.Stdout, strings.Join(prefixParts, "-"), log.LstdFlags|log.Lmsgprefix|log.LUTC|log.Lmicroseconds)

		if requestDigest != "" {
			logger.Printf("digest: %s", requestDigest)
			validateDigest = true
		}

		logger.Printf("method: %s, content length: %d\n", r.Method, r.ContentLength)

		// handle request body
		for {
			n, err := r.Body.Read(buffer)

			totalBytes += n
			logger.Printf("read %d bytes\n", n)
			if validateDigest {
				hasher.Write(buffer[:n])
			}

			if err != nil {
				if !errors.Is(err, io.EOF) {
					logger.Printf("error reading after %d bytes: %s\n", totalBytes, err.Error())
					responseStatus = http.StatusBadRequest
				}
				break
			}
		}

		if validateDigest {
			calculatedDigest := fmt.Sprintf("%x", hasher.Sum(nil))
			if requestDigest != calculatedDigest {
				logger.Printf("digest mismatch, request: %s vs calculated: %s", requestDigest, calculatedDigest)
				responseStatus = http.StatusBadRequest
			}
		}

		if sum > 0 {
			logger.Printf("read %d total bytes\n", totalBytes)
		}

		w.WriteHeader(responseStatus)
		_, _ = w.Write([]byte(""))
	}

	http.Serve(listener, handler)
}
