package main 

import (
	"fmt"
	"net"
)

type Address struct  {
	IP net.IP
	Port int
}

func (address *Address) String() string {
	return fmt.Sprintf("%s:%d", address.IP.String(), address.Port)
}