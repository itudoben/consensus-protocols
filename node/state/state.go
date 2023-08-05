package state

import "fmt"
import "net"

func init() {
	fmt.Println("itudoben.io/state package initialized")
}

type State struct {
	Ip   net.IP
	Role string
}
