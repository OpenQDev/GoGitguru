package ratelimit

import (
	"net"
	"net/http"
)

func getClientIP(r *http.Request) string {
	// X-Forwarded-For (XFF) HTTP header field is a de facto standard
	// for identifying the originating IP address of a client connecting
	// to a web server through an HTTP proxy or load balancer.
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		return xForwardedFor
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}
