package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var lastData []byte
var mu *sync.Mutex

func logErrorFf(server string, msg string, msgArgs ...any) {
	prefix := fmt.Sprintf("[%s][%s][ERROR]: ", time.Now().Format(time.RFC3339), server)
	log.Println(prefix + fmt.Sprintf(msg, msgArgs...))
}

func main() {
	var err error

	portForPrometheus := 0
	if portForPrometheus, err = strconv.Atoi(os.Getenv("PORT_FOR_PROMETHEUS")); err != nil || portForPrometheus <= 0 {
		log.Fatal("PORT_FOR_PROMETHEUS environment variable must be positive integer")
		// fatal
	}

	portForPush := 0
	if portForPush, err = strconv.Atoi(os.Getenv("PORT_FOR_PUSH_FROM_LOCALHOST")); err != nil || portForPush <= 0 {
		log.Fatal("PORT_FOR_PUSH_FROM_LOCALHOST environment variable must be positive integer")
		// fatal
	}

	log.Println(fmt.Sprintf("Port for Prometheus: %d", portForPrometheus))
	log.Println(fmt.Sprintf("Port for Push from Localhost: %d", portForPush))

	// ----

	lastData = make([]byte, 0)
	mu = new(sync.Mutex)

	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		defer wg.Done()
		serverForPrometheus(portForPrometheus)
	}()
	go func() {
		defer wg.Done()
		serverForPushFromLocalhost(portForPush)
	}()
	wg.Wait()
}

func serverForPrometheus(port int) {
	server := "prometheus"
	mux := http.NewServeMux()

	mux.Handle("GET /metrics", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		if len(lastData) == 0 {
			w.WriteHeader(http.StatusFailedDependency)
			return
		}

		n, err := w.Write(lastData)
		if err != nil {
			logErrorFf(server, "Error reading body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else if n != len(lastData) {
			logErrorFf(server, "Error writing body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}))

	log.Println(fmt.Sprintf("Server %s start...", server))
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		log.Fatal("Error starting "+server+" HTTP server", err)
		// fatal
	}
}

func serverForPushFromLocalhost(port int) {
	server := "push"
	mux := http.NewServeMux()

	mux.Handle("POST /", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		metricsData, err := io.ReadAll(r.Body)
		if err != nil {
			logErrorFf(server, "Error reading body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if len(metricsData) == 0 {
			logErrorFf(server, "Empty body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		lastData = metricsData
	}))

	log.Println(fmt.Sprintf("Server %s start...", server))
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		log.Fatal("Error starting "+server+" HTTP server", err)
		// fatal
	}
}
