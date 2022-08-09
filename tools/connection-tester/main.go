package main

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/digitalocean/doks-debug/tools/connection-tester/cmd"
)

func main() {
	ctx := context.Background()

	http.Handle("/metrics", promhttp.Handler())
	server := &http.Server{Addr: ":9100", Handler: nil}
	go server.ListenAndServe()
	defer server.Shutdown(ctx)

	cmd.Execute(ctx)
}
