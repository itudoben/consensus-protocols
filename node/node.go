package main

import (
	"fmt"
	"io"
	"itudoben.io/state"
	"log"
	"net"
	"net/http"
	"os"
)

// Structure to implement interface http.ResponseWriter
type Dummy struct {
	w io.Writer
}

// This is how to implement an interface by just providing methods of the interface one wants
// to implement and attached to the struct passed
// just after keyword func i.e. d Dummy.
func (d Dummy) Header() http.Header          { return nil }
func (d Dummy) Write(bb []byte) (int, error) { return d.w.Write(bb) }
func (d Dummy) WriteHeader(statusCode int)   {}

var p = 8000

var stat = new(state.State)

func main() {
	setState(stat, os.Stdout)

	http.HandleFunc("/", defaultHandler) // each request calls handler
	http.HandleFunc("/status", status)   // each request calls handler

	log.Print()
	log.Fatal(http.ListenAndServe(":"+fmt.Sprint(p), nil))
}

// handler echoes the Path component of the request URL r.
func defaultHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "[DEFAULT] URL Path: %s\n", req.URL.Path)
}

func status(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "[STATUS] Node with IP %+v:%d\n", stat, p)
}

func setState(stat *state.State, w io.Writer) {
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
				stat.Ip = ipnet.IP.String()
				stat.Role = "subsidiary"
			}
		}
	}
}
