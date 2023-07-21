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
var pUdp = 8972

var stat = new(state.State)

func httpServer() error {
	setState(stat, os.Stdout)

	http.HandleFunc("/", defaultHandler) // each request calls handler
	http.HandleFunc("/status", status)   // each request calls handler

	thisIp, err := GetLocalIP(os.Stdout)
	if err != nil {
		return err
	}

	fmt.Printf("%s waiting for HTTP commands on port %d...\n", thisIp, p)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", p), nil))
	return nil
}

func listerBroadcast() error {
	thisIp, err := GetLocalIP(os.Stdout)
	if err != nil {
		return err
	}

	pc, err := net.ListenPacket("udp4", fmt.Sprintf(":%d", pUdp))
	defer pc.Close()
	if err != nil {
		return err
	}

	buf := make([]byte, 1024)
loop:
	for {
		fmt.Printf("%s waiting for broadcast commands on port %d...\n", thisIp, pUdp)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			return err
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
			fmt.Printf("%s shutdown\n", thisIp)
			break loop
		case "i":
			if err != nil {
				return err
			}

			fmt.Printf("%s received %q sent by %s\n", thisIp, c, addr.String())
		default:
			fmt.Printf("%s received unknown command %q sent by %s\n", thisIp, c, addr.String())
		}
	}
	return nil
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

func setState(stat *state.State, w io.Writer) error {
	ip, err := GetLocalIP(w)
	if err != nil {
		return err
	}

	stat.Ip = ip
	stat.Role = "subsidiary"
	return nil
}

func main() {
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		err := httpServer()
		defer wg.Done()
		if err != nil {
			panic(err)
		}
	}()

	wg.Add(1)
	go func() {
		err := listerBroadcast()
		defer wg.Done()
		if err != nil {
			panic(err)
		}
	}()

	wg.Wait()
}
