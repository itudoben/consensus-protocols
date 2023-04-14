package state

import "fmt"

func init() {
	fmt.Println("itudoben.io/state package initialized")
}

type State struct {
	Ip   string
	Role string
}
