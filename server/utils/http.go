package utils

import "net/http"

// Simple way to get client ip address from http.Request
// TODO: What if client set this header and we forget to setup nginx to write this header
func GetIP(r *http.Request) string {
	ips, ok := r.Header["X-Real-IP"]
	if ok && len(ips) > 0 {
		return ips[0]
	}
	return r.RemoteAddr
}
