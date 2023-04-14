package main

import (
	"fmt"
	"log"
	"net"
	"io"
	"os"
	"net/http"
)

type Dummy struct {
    w io.Writer
}

func (d Dummy) Header() http.Header { return nil}
func (d Dummy) Write(bb []byte) (int, error) { return d.w.Write(bb)}
func (d Dummy) WriteHeader(statusCode int) {}

var p = 8000

func main() {
    status(Dummy{os.Stdout}, nil)

	http.HandleFunc("/", handler)      // each request calls handler
	http.HandleFunc("/status", status) // each request calls handler


	// TODO: pass the state in the status handler to set the state of the server and edit it.
// 	state := new(state.State)


	log.Print()
	log.Fatal(http.ListenAndServe(":" + fmt.Sprint(p), nil))
}

// handler echoes the Path component of the request URL r.
func handler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "URL Path: %q\n", req.URL.Path)
}

func status(w http.ResponseWriter, req *http.Request) {
	// get list of available addresses
	addr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	for _, addr := range addr {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			// check if IPv4 or IPv6 is not nil
			if ipnet.IP.To4() != nil || ipnet.IP.To16 != nil {
				// print available addresses
                fmt.Printf("[STATUS] %s received request\n", ipnet.IP.String())
				fmt.Fprintf(w, "[STATUS] Node with IP %s:%d\n", ipnet.IP.String(), p)
			}
		}
	}
}
