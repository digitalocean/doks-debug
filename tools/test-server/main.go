package main

import (
	"context"
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
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpConnectionsHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: "testserver",
		Name:      "http_connections_duration_s",
		Help:      "HTTP connection establishment duration in seconds",
		//     >>> for i in range(0, 20):
		// ...   print("%.2f" % (0.01 * 1.9 ** i))
		// ...
		// 0.01
		// 0.02
		// 0.04
		// 0.07
		// 0.13
		// 0.25
		// 0.47
		// 0.89
		// 1.70
		// 3.23
		// 6.13
		// 11.65
		// 22.13
		// 42.05
		// 79.90
		// 151.81
		// 288.44
		// 548.04
		// 1041.27
		// 1978.42
		Buckets: prometheus.ExponentialBuckets(0.01, 1.9, 20),
	},
		[]string{
			"method",
			"status_code",
			"result_type",
			"content_length",
			"source_region",
			"region",
		},
	)
)

func main() {
	ctx := context.Background()
	http.Handle("/metrics", promhttp.Handler())
	server := &http.Server{Addr: ":9100", Handler: nil}
	go server.ListenAndServe()
	defer server.Shutdown(ctx)

	listenHTTP()
}

func listenHTTP() {
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

			requestDigest = r.Header.Get("digest")
			sourceRegion  = r.Header.Get("region")
			region        = os.Getenv("REGION")

			validateDigest bool

			buffer         = make([]byte, 1024+1024)
			totalBytes     int64
			hasher         = sha256.New()
			responseStatus = http.StatusOK
			startTime      = time.Now()
			resultType     = "success"
		)

		if sourceRegion == "" {
			sourceRegion = "unknown"
		}

		if region == "" {
			region = "unknown"
		}

		defer func() {
			labelValues := []string{
				r.Method,
				fmt.Sprintf("%d", responseStatus),
				resultType,
				fmt.Sprintf("%d", r.ContentLength),
				sourceRegion,
				region,
			}
			duration := time.Now().Sub(startTime).Seconds()
			httpConnectionsHistogram.WithLabelValues(labelValues...).Observe(duration)
		}()

		if requestUuid != "" {
			prefixParts = append(prefixParts, requestUuid)
		}
		logger := log.New(os.Stdout, strings.Join(prefixParts, "-")+" ", log.LstdFlags|log.Lmsgprefix|log.LUTC|log.Lmicroseconds)

		if requestDigest != "" {
			logger.Printf("digest: %s", requestDigest)
			validateDigest = true
		}

		logger.Printf("method: %s, content length: %d\n", r.Method, r.ContentLength)

		// handle request body
		for {
			n, err := r.Body.Read(buffer)

			totalBytes += int64(n)
			//logger.Printf("read %d bytes\n", n)
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

		if totalBytes != r.ContentLength {
			logger.Printf("content length mismatch, request: %d vs calculated: %d", r.ContentLength, totalBytes)
			resultType = "content length mismatch"
			responseStatus = http.StatusBadRequest
		} else if validateDigest {
			calculatedDigest := fmt.Sprintf("%x", hasher.Sum(nil))
			if requestDigest != calculatedDigest {
				logger.Printf("digest mismatch, request: %s vs calculated: %s", requestDigest, calculatedDigest)
				resultType = "digest mismatch"
				responseStatus = http.StatusBadRequest
			}
			logger.Printf("request and calculated digest match")
		}

		if totalBytes > 0 {
			logger.Printf("read %d total bytes\n", totalBytes)
		}

		w.WriteHeader(responseStatus)
		_, _ = w.Write([]byte(""))
	}

	http.Serve(listener, handler)
}
