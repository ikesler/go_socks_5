package main

import (
	"context"
	"github.com/things-go/go-socks5/statute"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	socks5 "github.com/things-go/go-socks5"
)

type MyRewriter struct {
	Hosts        []string
	UpstreamHost net.IP
	UpstreamPort int
	Ips          **[]net.IP
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
			IP:   rw.UpstreamHost,
			Port: rw.UpstreamPort,
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
	println("Starting application")
	hostsEnv := os.Getenv("ROUTE_TO_UPSTREAM_HOSTS")
	upstreamIpEnv := os.Getenv("UPSTREAM_IP")
	upstreamPortEnv := os.Getenv("UPSTREAM_PORT")

	var rewriter MyRewriter

	if hostsEnv != "" && upstreamIpEnv != "" && upstreamPortEnv != "" {
		println("Running resolver routine")
		hosts := strings.Split(hostsEnv, ",")
		upstreamIp := net.ParseIP(upstreamIpEnv)
		var upstreamPort int64
		upstreamPort, _ = strconv.ParseInt(upstreamPortEnv, 10, 32)

		ips := []net.IP{}
		ipsRef := &ips

		rewriter = MyRewriter{
			Hosts:        hosts,
			Ips:          &ipsRef,
			UpstreamHost: upstreamIp,
			UpstreamPort: int(upstreamPort),
		}
		go rewriter.resolve()
	}

	// Create a SOCKS5 server
	server := socks5.NewServer(
		socks5.WithLogger(socks5.NewLogger(log.New(os.Stdout, "socks5: ", log.LstdFlags))),
		socks5.WithRewriter(rewriter),
	)

	println("Running proxy routine")
	if err := server.ListenAndServe("tcp", ":10800"); err != nil {
		panic(err)
	}
}
