package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)      // each request calls handler
	http.HandleFunc("/status", status) // each request calls handler



	// TODO: pass the state in the status handler to set the state of the server and edit it.
	state := new(state.State)




	log.Fatal(http.ListenAndServe(":8000", nil))
}

// handler echoes the Path component of the request URL r.
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL Path: %q\n", r.URL.Path)
}

func status(w http.ResponseWriter, r *http.Request) {
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
				fmt.Printf("Node with IP %s handles request\n", ipnet.IP.String())
				fmt.Fprintf(w, "Hello from node with IP %s\n", ipnet.IP.String())
			}
		}
	}
}
