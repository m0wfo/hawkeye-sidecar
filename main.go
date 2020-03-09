package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/hpcloud/tail"
)

func main() {
	fmt.Println("Sidecar started")

	conn, err := net.Dial("tcp", "localhost:8080")
	handleErr(err)

	defer conn.Close()

	req, err := http.NewRequest("GET", "http://localhost:8080", nil)
	handleErr(err)

	req.Header.Add("Connection", "Upgrade")
	req.Header.Add("Upgrade", "hawkeye/1.0.0alpha1")

	handleErr(err)

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	err = req.Write(writer)
	handleErr(err)
	writer.Flush()

	resp, err := http.ReadResponse(reader, &http.Request{Method: "GET"})
	handleErr(err)

	fmt.Println(resp.Status)

	controlMessage := make(map[string]string)
	encoder := json.NewEncoder(writer)

	controlMessage["__hawkeye_filename"] = "/var/log/foo.test"

	err = encoder.Encode(controlMessage)
	handleErr(err)
	writer.Flush()

	ok, err := reader.ReadString('\n')
	handleErr(err)
	if ok == "OK\n" {
		fmt.Println("came out ok")
	}

	t, err := tail.TailFile("/Users/chris/Desktop/test.log", tail.Config{Follow: true})
	for line := range t.Lines {
		writer.WriteString(line.Text)
		writer.Flush()
		fmt.Println(line.Text)
	}
}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
