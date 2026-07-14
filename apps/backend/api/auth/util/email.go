package auth_util

import (
	"context"
	"net"
	"strings"
	"time"
)

func ValidateEmailMX(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	domain := parts[1]

	resolver := net.Resolver{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mx, err := resolver.LookupMX(ctx, domain)
	if err == nil && len(mx) > 0 {
		return true
	}

	ips, err := resolver.LookupIPAddr(ctx, domain)
	if err == nil && len(ips) > 0 {
		return true
	}

	return false
}
