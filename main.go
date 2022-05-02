package main

import (
	socks5 "github.com/things-go/go-socks5"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	println("Starting application")
	hostsEnv := os.Getenv("ROUTE_TO_UPSTREAM_HOSTS")
	upstreamEnv := os.Getenv("UPSTREAM")

	if hostsEnv != "" && upstreamEnv != "" {
		println("Running resolver routine")
		hosts := strings.Split(hostsEnv, ",")

		ips := []net.IP{}
		ipsRef := &ips

		upstream := socks5.UpstreamProxy{
			Hosts:            hosts,
			Ips:              &ipsRef,
			UpstreamEndpoint: upstreamEnv,
		}
		go upstream.Resolve()

		// Create a SOCKS5 server
		server := socks5.NewServer(
			socks5.WithLogger(socks5.NewLogger(log.New(os.Stdout, "socks5: ", log.LstdFlags))),
			socks5.WithUpstream(upstream),
		)
		println("Running proxy routine (with upstream)")
		if err := server.ListenAndServe("tcp", ":10800"); err != nil {
			panic(err)
		}
	} else {
		// Create a SOCKS5 server
		server := socks5.NewServer(
			socks5.WithLogger(socks5.NewLogger(log.New(os.Stdout, "socks5: ", log.LstdFlags))),
		)
		println("Running proxy routine (without upstream)")
		if err := server.ListenAndServe("tcp", ":10800"); err != nil {
			panic(err)
		}
	}
}
