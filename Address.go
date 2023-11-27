package main 

import (
	"fmt"
	"net"
)

type Address struct  {
	IP net.IP
	Port int
}

func NewAddress(port int) Address {
	return Address{
		IP: net.IPv4(127, 0, 0, 1),
		Port: port,
	}
}

func (address *Address) String() string {
	return fmt.Sprintf("%s:%d", address.IP.String(), address.Port)
}