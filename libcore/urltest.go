package libcore

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/v2fly/v2ray-core/v5/common/net"
	"github.com/v2fly/v2ray-core/v5/common/session"
)

func UrlTestV2ray(instance *V2RayInstance, inbound string, link string, timeout int32) (int32, error) {
	// connTestUrl, err := url.Parse(link)
	// address := net.ParseAddress(connTestUrl.Hostname())

	transport := &http.Transport{
		TLSHandshakeTimeout: time.Duration(timeout) * time.Millisecond,
		DisableKeepAlives:   true,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// I believe the server...

			// if address.Family().IsDomain() {
			// 	ips, err := dc.LookupIP(&dns.MatsuriDomainStringEx{
			// 		Domain:  address.Domain(),
			// 		Network: "ip4",
			// 	})
			// 	if err != nil {
			// 		return nil, err
			// 	} else if len(ips) == 0 {
			// 		return nil, newError("no ip")
			// 	}
			// 	addr2, _ := net.ParseDestination(addr)
			// 	addr = fmt.Sprintf("%s:%s", ips[0].String(), addr2.Port.String())
			// }

			dest, err := net.ParseDestination(fmt.Sprintf("%s:%s", network, addr))
			if err != nil {
				return nil, err
			}
			if inbound != "" {
				ctx = session.ContextWithInbound(ctx, &session.Inbound{Tag: inbound})
			}
			return instance.dialContext(ctx, dest)
		},
	}
	req, err := http.NewRequestWithContext(context.Background(), "GET", link, nil)
	req.Header.Set("User-Agent", fmt.Sprintf("curl/7.%d.%d", rand.Int()%54, rand.Int()%2))
	if err != nil {
		return 0, err
	}
	start := time.Now()
	resp, err := (&http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeout) * time.Millisecond,
	}).Do(req)
	if err == nil && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unexcpted response status: %d", resp.StatusCode)
	}
	if err != nil {
		return 0, err
	}
	return int32(time.Since(start).Milliseconds()), nil
}
