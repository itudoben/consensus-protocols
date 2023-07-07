package main

import (
	"fmt"
	"io"
	"itudoben.io/state"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
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
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		httpServer()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		listerBroadcast()
		wg.Done()
	}()

	wg.Wait()
}

func httpServer() {
	setState(stat, os.Stdout)

	http.HandleFunc("/", defaultHandler) // each request calls handler
	http.HandleFunc("/status", status)   // each request calls handler

    fmt.Printf("%s waiting for broadcast commands ...\n", addrLoc)

	log.Fatal(http.ListenAndServe(":"+fmt.Sprint(p), nil))
}

func listerBroadcast() {
	thisIp, errs := GetLocalIP(os.Stdout)

	pc, err := net.ListenPacket("udp4", ":8872")
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	buf := make([]byte, 1024)
	addrLoc := thisIp
loop:
	for {
		fmt.Printf("%s waiting for broadcast commands ...\n", addrLoc)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			panic(err)
			break
		}

		c := ""
		var x uint8 = 0
		// Clean up what was received to remove an added character at the end of the buffer
		if buf[n] == x {
			c = string(buf[:n-1])
		} else {
			c = string(buf[:n])
		}

		switch c {
		case "q":
			fmt.Printf("%s shutdown\n", addrLoc)
			break loop
		case "i":
			if errs != nil {
				panic(err)
			}

			fmt.Printf("%s received %q sent by %s\n", thisIp, c, addr.String())
		default:
			fmt.Printf("%s received unknown command %q sent by %s\n", addrLoc, c, addr.String())
		}
	}
}

// handler echoes the Path component of the request URL r.
func defaultHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "[DEFAULT] URL Path: %s\n", req.URL.Path)
}

func status(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(os.Stdout, "[STATUS] Node with IP %+v:%d\n", stat, p)
	fmt.Fprintf(w, "[STATUS] Node with IP %+v:%d\n", stat, p)
}

func GetLocalIP(w io.Writer) (ip string, err error) {
	// Get preferred outbound ip of this machine
	conn, err := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

func setState(stat *state.State, w io.Writer) {
	ip, err := GetLocalIP(w)
	if err != nil {
		panic(err)
	}

	stat.Ip = ip
	stat.Role = "subsidiary"
}
