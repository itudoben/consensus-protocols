package main

import (
	"fmt"
	// 	netaddr "gopkg.in/netaddr.v1"
	"io"
	"itudoben.io/state"
	"log"
	"math/rand"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sync"
	"time"
	netstate "v.io/x/lib/netstate"
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

type Node struct {
	portHttp int
	portUDP  int
	ip       *net.IP
}

func NewNode(portHttp int, portUDP int, ip *net.IP) *Node {
	return &Node{portHttp: portHttp, portUDP: portUDP, ip: ip}
}

var stat = new(state.State)

func httpServer(node *Node) error {
	setState(stat, os.Stdout)

	http.HandleFunc("/", defaultHandler)              // each request calls handler
	http.Handle("/status", &countHandler{node: node}) // each request calls handler

	// Get the local machine's hostname
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Printf("%s (%s) waiting for HTTP commands on port %d...\n", hostname, node.ip.String(), node.portHttp)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", node.portHttp), nil))
	return nil
}

// The broadcast allows to give commands to the server config via the channel
var sem = make(chan string, 0)

func listerBroadcast(node *Node) error {
	pc, err := net.ListenPacket("udp4", fmt.Sprintf(":%d", node.portUDP))
	defer pc.Close()
	if err != nil {
		return err
	}

	buf := make([]byte, 1024)
	ifcs, _ := netstate.GetAllInterfaces()
	for i := 0; i < len(ifcs); i++ {
		addrs := ifcs[i].Addrs()
		for j := 0; j < len(addrs); j++ {
			fmt.Printf("%d %d net %s, str %s \n", i, j, addrs[j].Network(), addrs[j].String())
		}
	}

loop:
	for {
		// Get the local machine's hostname
		hostname, err := os.Hostname()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		fmt.Printf("%s (%s) waiting for broadcast from ''%s' commands on port %d...\n", hostname, node.ip.String(),
			"netaddr.BroadcastAddr(net.Addr.Network()).String()", node.portUDP)
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
			fmt.Printf("%s shutdown\n", node.ip.String())
			break loop
		case "i":
			if err != nil {
				return err
			}

			fmt.Printf("[BROADCAST] %s (%s) received %q sent by %s\n", hostname, node.ip.String(), c, addr.String())
		default:
			fmt.Printf("[BROADCAST] %s (%s) received unknown command %q sent by %s\n", hostname, node.ip.String(), c, addr.String())
		}
	}
	return nil
}

// handler echoes the Path component of the request URL r.
func defaultHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "[HTTP] URL Path: %s\n", req.URL.Path)
}

type countHandler struct {
	node *Node
}

func (h *countHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(os.Stdout, "[STATUS] Node with IP %+v:%d\n", stat, h.node.portHttp)
	fmt.Fprintf(w, "[STATUS] Node with IP %+v:%d\n", stat, h.node.portHttp)
}

func GetLocalIP(w io.Writer) (ip *net.IP, err error) {
	// Get preferred outbound ip of this machine
	conn, err := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return &localAddr.IP, nil
}

func setState(stat *state.State, w io.Writer) error {
	ip, err := GetLocalIP(w)
	if err != nil {
		return err
	}

	stat.Ip = *ip
	stat.Role = "follower" // leader as in Raft consensus algorithm leader election
	return nil
}

// 150 to 300 ms randomized timeout
var electionTimeout = randomValue(150, 300)

func heartBeat() {
	for i := 0; i < 10; i++ {
		fmt.Printf("Hearbeat in %s ...\n", electionTimeout)
		time.Sleep(electionTimeout)
		callFunction()
	}

	// return &error.Error("Fail to get the heartBeat")
}

func randomValue(min int, max int) time.Duration {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate a random integer between 150 and 300
	return time.Duration(rand.Intn(min)+max-min) * time.Millisecond // Generates a value between 0 and 150, then adds 150
}

func callFunction() {
	fmt.Println("Function called now")
}

func main() {
	// First send a request to join a cluster ID
	//     broadcastAddress, err := netaddr.Broadcast.String()
	//
	//     	conn, err := net.Dial("udp", broadcastAddress)
	//     	if err != nil {
	//     		fmt.Println("Error creating UDP connection:", err)
	//     		return
	//     	}
	//     	defer conn.Close()
	//
	//     	message := []byte("This is a broadcast message!")
	//     	_, err = conn.Write(message)
	//     	if err != nil {
	//     		fmt.Println("Error sending UDP broadcast:", err)
	//     		return
	//     	}
	//
	//     	fmt.Println("Broadcast message sent successfully.")

	thisIp, err := GetLocalIP(os.Stdout)
	if err != nil {
		panic(err)
	}

	thisNode := NewNode(8000, 8972, thisIp)
	wg := &sync.WaitGroup{}

	go func() { log.Println(http.ListenAndServe("localhost:6060", nil)) }()

	wg.Add(1)
	go func() {
		err := httpServer(thisNode)
		defer wg.Done()
		if err != nil {
			panic(err)
		}
	}()

	wg.Add(1)
	go func() {
		err := listerBroadcast(thisNode)
		defer wg.Done()
		if err != nil {
			panic(err)
		}
	}()

	wg.Add(1)
	go func() {
		heartBeat()
	}()

	wg.Wait()
}
