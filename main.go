package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hpcloud/tail"
	"github.com/tuplestream/hawkeye-client"
)

var authToken = os.Getenv("TUPLESTREAM_AUTH_TOKEN")

func main() {
	fmt.Println("Sidecar started")
	pattern := "/var/log/containers/*.log"
	matches, err := filepath.Glob(pattern)
	handleErr(err)

	for _, path := range matches {
		if !strings.Contains(path, "hawkeye-sidecar") {
			go tailFile(path)
		}
	}

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	})

	log.Fatal(http.ListenAndServe(":8889", nil))
}

func tailFile(filename string) {
	log.Print("Starting to tail " + filename)
	conn, writer := hawkeye.InitiateConnection(filename, authToken)
	for {
		if writer != nil {
			break
		}
		log.Print("retrying connection in 5 seconds")
		time.Sleep(5 * time.Second)
		conn, writer = hawkeye.InitiateConnection(filename, authToken)
	}
	defer conn.Close()
	shouldRetry := false
	t, err := tail.TailFile(filename, tail.Config{Follow: true})
	handleErr(err)
	for line := range t.Lines {

		written, err := writer.WriteString(line.Text + "\n")
		if written < len(line.Text) || err != nil {
			log.Print("Connection closed")
			shouldRetry = true
			break
		}
		writer.Flush()
	}
	if shouldRetry {
		log.Print("Retrying connection for " + filename)
		tailFile(filename)
	}
}

func handleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
