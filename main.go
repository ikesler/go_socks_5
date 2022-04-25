package main

import (
	"context"
	"github.com/things-go/go-socks5/statute"
	"log"
	"net"
	"os"
	"strings"
	"time"

	socks5 "github.com/things-go/go-socks5"
)

type MyRewriter struct {
	Hosts    []string
	Upstream net.IP
	Ips      **[]net.IP
}

func (rw MyRewriter) shouldRoute(ip net.IP) bool {
	if rw.Ips == nil {
		return false
	}
	for _, v := range **rw.Ips {
		if ip.Equal(v) {
			return true
		}
	}

	return false
}

func (rw MyRewriter) Rewrite(ctx context.Context, request *socks5.Request) (context.Context, *statute.AddrSpec) {
	if rw.shouldRoute(request.DestAddr.IP) {
		result := &statute.AddrSpec{
			IP: rw.Upstream,
		}
		return ctx, result
	}

	return ctx, request.DestAddr
}

func (rw MyRewriter) resolve() {
	for _, host := range rw.Hosts {
		t, _ := net.LookupIP(host)
		newIps := append(**rw.Ips, t...)
		*rw.Ips = &newIps
	}

	time.Sleep(10 * 1000)
}

func main() {
	hostsEnv := os.Getenv("ROUTE_TO_UPSTREAM_HOSTS")
	upstreamEnv := os.Getenv("UPSTREAM_IP")

	hosts := strings.Split(hostsEnv, ",")
	upstream := net.ParseIP(upstreamEnv)

	ips := []net.IP{}
	ipsRef := &ips

	rewriter := MyRewriter{
		Hosts:    hosts,
		Ips:      &ipsRef,
		Upstream: upstream,
	}
	go rewriter.resolve()
	// Create a SOCKS5 server
	server := socks5.NewServer(
		socks5.WithLogger(socks5.NewLogger(log.New(os.Stdout, "socks5: ", log.LstdFlags))),
		socks5.WithRewriter(rewriter),
	)

	if err := server.ListenAndServe("tcp", ":10800"); err != nil {
		panic(err)
	}
}
