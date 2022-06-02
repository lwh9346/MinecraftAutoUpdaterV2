package main

import (
	"flag"
	"mau2/proxy"
)

func main() {
	var (
		listen = flag.String("l", "", "listen address")
		target = flag.String("t", "", "target address")
	)
	flag.Parse()
	proxy.SetUp(*listen, *target, true)
}
