package check_response

import (
    "context"
    "net"
    "strings"

    "github.com/coredns/coredns/plugin"
    "github.com/coredns/coredns/request"
    "github.com/miekg/dns"
)

type CheckResponse struct {
    Next plugin.Handler
}

// ServeDNS implements the plugin.Handler interface.
func (c CheckResponse) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
    state := request.Request{W: w, Req: r}

    // Call the next plugin in the chain and get the result
    rc, err := plugin.NextOrFailure(c.Name(), c.Next, ctx, w, r)
    if err != nil {
        return rc, err
    }

    // Check if the response is of type A (IPv4)
    for _, answer := range w.(*dns.Msg).Answer {
        if a, ok := answer.(*dns.A); ok {
            ip := a.A.String()

            // Check if the IP is in the 10.10.*.* range
            if strings.HasPrefix(ip, "10.10.") {
                // Log and switch to the next resolver if the IP matches 10.10.*.*
                plugin.Logger(c).Infof("IP %s matches 10.10.*.*, forwarding to next resolver", ip)
                state.Do() // Re-run the DNS query using the next resolver
            }
        }
    }

    return rc, nil
}

// Name returns the plugin name.
func (c CheckResponse) Name() string { return "check_response" }
