package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var portExporter int
var portProxyPush int
var intervalSeconds int

func logErrorFf(tickId int64, msg string, msgArgs ...any) {
	prefix := fmt.Sprintf("[%s][%d][ERROR]: ", time.Now().Format(time.RFC3339), tickId)
	log.Println(prefix + fmt.Sprintf(msg, msgArgs...))
}

func main() {
	var err error

	portExporter = 0
	if portExporter, err = strconv.Atoi(os.Getenv("PORT_EXPORTER")); err != nil || portExporter <= 0 {
		log.Fatal("PORT_EXPORTER environment variable must be positive integer")
		// fatal
	}

	portProxyPush = 0
	if portProxyPush, err = strconv.Atoi(os.Getenv("PORT_PROXY_PUSH")); err != nil || portProxyPush <= 0 {
		log.Fatal("PORT_PROXY_PUSH environment variable must be positive integer")
		// fatal
	}

	if intervalSeconds, err = strconv.Atoi(os.Getenv("INTERVAL_SECONDS")); err != nil || intervalSeconds <= 0 {
		log.Fatal("INTERVAL_SECONDS environment variable must be positive integer")
		// fatal
	}

	log.Println(fmt.Sprintf("Port Exporter: %d", portExporter))
	log.Println(fmt.Sprintf("Port Proxy Push: %d", portProxyPush))
	log.Println(fmt.Sprintf("Interval Seconds: %d", intervalSeconds))

	// ----

	log.Println("Go (first)...")
	handleTick(time.Now())

	log.Println("Go (further)...")
	for tick := range time.Tick(time.Duration(intervalSeconds) * time.Second) {
		handleTick(tick)
	}
}

func handleTick(tick time.Time) {
	tickId := tick.Unix()

	metricsData := pullMetrics(tickId)
	if len(metricsData) == 0 {
		logErrorFf(tickId, "Failed to PULL metricsData: no metricsData")
		return
	}

	err := pushToProxy(tickId, metricsData)
	if err != nil {
		logErrorFf(tickId, "Failed to PUSH metricsData: %v", err)
		return
	}
}

func pullMetrics(tickId int64) []byte {
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/metrics", portExporter))

	if err != nil {
		logErrorFf(tickId, "Failed to fetch metricsData: %v", err)
		return nil
	} else if resp.StatusCode != 200 {
		logErrorFf(tickId, "Response is not ok: %s", resp.Status)
		return nil
	}

	defer func(resp *http.Response) {
		err := resp.Body.Close()
		if err != nil {
			logErrorFf(tickId, "Failed to close response body: %v", err)
			return
		}
	}(resp)

	metricsData, err := io.ReadAll(resp.Body)
	if err != nil {
		logErrorFf(tickId, "Failed to read response body: %v", err)
		return nil
	} else if len(metricsData) == 0 {
		logErrorFf(tickId, "Empty response body %v", metricsData)
		return nil
	}

	return metricsData
}

func pushToProxy(tickId int64, metricsData []byte) error {
	resp, err := http.Post(fmt.Sprintf("http://127.0.0.1:%d/", portProxyPush), "text/html", bytes.NewBuffer(metricsData))

	if err != nil {
		logErrorFf(tickId, "Failed to push metrics to proxy: %v", err)
		return err
	} else if resp.StatusCode != 200 {
		logErrorFf(tickId, "Response is not ok: %s", resp.Status)
		return err
	}

	defer func(resp *http.Response) {
		err := resp.Body.Close()
		if err != nil {
			logErrorFf(tickId, "Failed to close response body: %v", err)
			return
		}
	}(resp)

	return nil
}
