package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hpcloud/tail"
)

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
	conn, writer := initiateConnection(filename)
	for {
		if writer != nil {
			break
		}
		log.Print("retrying connection in 5 seconds")
		time.Sleep(5 * time.Second)
		conn, writer = initiateConnection(filename)
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

func initiateConnection(filename string) (net.Conn, *bufio.Writer) {
	hawkeyeTarget := getHawkeyeTarget()
	conn, err := net.Dial("tcp", hawkeyeTarget.Host)

	if err != nil {
		log.Print("error connecting to " + hawkeyeTarget.Host)
		return nil, nil
	}

	req, err := http.NewRequest("GET", getHawkeyeTarget().String(), nil)
	handleErr(err)

	req.Header.Add("Connection", "Upgrade")
	req.Header.Add("Upgrade", "hawkeye/1.0.0alpha1")
	req.Header.Add("User-Agent", "hawkeye/client-go1.0.0alpha1")

	handleErr(err)

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	err = req.Write(writer)
	handleErr(err)
	writer.Flush()

	resp, err := http.ReadResponse(reader, req)
	if resp.StatusCode != 101 {
		log.Fatal("Couldn't upgrade HTTP connection, closing. Got status: " + resp.Status)
	}
	handleErr(err)

	fmt.Println(resp.Status)

	controlMessage := make(map[string]string)
	encoder := json.NewEncoder(writer)

	controlMessage["__hawkeye_filename"] = filename

	err = encoder.Encode(controlMessage)
	handleErr(err)
	writer.Flush()

	ok, err := reader.ReadString('\n')
	handleErr(err)
	if ok == "OK\n" {
		log.Print("handshake successful")
	}
	return conn, writer
}

func getHawkeyeTarget() *url.URL {
	rawURL := os.Getenv("HAWKEYE_TARGET")
	if rawURL == "" {
		log.Fatal("no HAWKEYE_TARGET host specified")
	}
	u, e := url.Parse(rawURL)
	handleErr(e)
	return u
}

func handleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
