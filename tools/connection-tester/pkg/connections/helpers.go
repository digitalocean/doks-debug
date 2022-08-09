package connections

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	attempt  uint64
	hostname string
	region   string

	connectionsHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: "connectiontester",
		Name:      "tls_connections_duration_s",
		Help:      "TLS connection establishment duration in seconds",
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
		Buckets: prometheus.ExponentialBuckets(0.01, 1.9, 10),
	},
		[]string{"connection_type", "status", "error", "hostname", "region"},
	)
)

type regionKey struct{}

func init() {
	hostname = os.Getenv("HOSTNAME")
	region = os.Getenv("REGION")
}

// TryLoops creates the specified number of "loops" aka goroutines and runs
// Connector.Connect repeatedly in each one.
func TryLoops(ctx context.Context, numLoops int, c Connector) {
	wg := sync.WaitGroup{}

	for i := 0; i < numLoops; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			tryLoop(ctx, i, c)
		}()
	}

	wg.Wait()
}

func tryLoop(ctx context.Context, n int, c Connector) {
	connType := c.Type()
	logger := log.New(os.Stdout, fmt.Sprintf("%s %d ", connType, n), log.LstdFlags|log.Lmsgprefix|log.LUTC|log.Lmicroseconds)
	defer func() {
		logger.Printf("ending %s loop %d\n", connType, n)
	}()
	logger.Printf("starting %s loop %d\n", connType, n)

	for {
		thisAttempt := atomic.AddUint64(&attempt, 1)
		startTime := time.Now()
		ctx := context.WithValue(ctx, regionKey{}, region)
		err := c.Connect(ctx)
		duration := time.Now().Sub(startTime).Seconds()
		labelValues := []string{connType, "success", "", hostname, region}
		if err != nil {
			labelValues[0] = "fail"
			labelValues[1] = errorLabelValue(err.Error())
			logger.Printf("%d FAIL: %s", thisAttempt, err)
		}
		connectionsHistogram.WithLabelValues(labelValues...).Observe(duration)
		if ctx.Err() != nil {
			return
		}
	}
}

func errorLabelValue(s string) string {
	const connectionResetByPeer = "connection reset by peer"

	if !strings.Contains(s, connectionResetByPeer) {
		return s
	}

	return "conn.Handshake failed: <host1> -> <host2>: read: connection reset by peer"
}
