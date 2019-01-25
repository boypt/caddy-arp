package arp

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/jstotz/arp"
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

type Arp struct {
	Next httpserver.Handler
}

// Init initializes the plugin
func init() {
	caddy.RegisterPlugin("arp", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {

	// refres every 1s
	arp.AutoRefresh(time.Second)

	// Add middleware
	cfg := httpserver.GetConfig(c)
	cfg.AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		return &Arp{
			Next: next,
		}
	})

	return nil
}

func (a Arp) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	a.lookupArp(w, r)
	return a.Next.ServeHTTP(w, r)
}

func (a Arp) lookupArp(w http.ResponseWriter, r *http.Request) {
	var clientMac string
	clientIP, _ := getClientIP(r)

	if len(clientIP) > 0 {
		// ensure is ipv4 addr
		netIP := net.ParseIP(clientIP)
		if netIP != nil {
			ip4 := netIP.To4()
			if ip4 != nil {
				clientMac = arp.Search(ip4.String())
			}
		}
	}

	replacer := newReplacer(r)
	replacer.Set("client_mac", clientMac)
	if rr, ok := w.(*httpserver.ResponseRecorder); ok {
		rr.Replacer = replacer
	}
}

func getClientIP(r *http.Request) (string, error) {

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		if serr, ok := err.(*net.AddrError); ok && serr.Err == "missing port in address" { // It's not critical try parse
			ip = r.RemoteAddr
		} else {
			log.Printf("Error when SplitHostPort: %v", serr.Err)
			return "", err
		}
	}

	return ip, nil
}

func newReplacer(r *http.Request) httpserver.Replacer {
	return httpserver.NewReplacer(r, nil, "")
}
