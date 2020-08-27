package closestRR

import (
	"github.com/caddyserver/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

// init registers this plugin.
func init() { plugin.Register("closestRR", setup) }

func setup(c *caddy.Controller) error {
	t, err := parse(c)
	if err != nil {
		return plugin.Error("closestRR", err)
	}
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		t.Next = next
		return t
	})

	return nil
}

func parse(c *caddy.Controller) (*closestRR, error) {
	found := false
	e := &closestRR{}

	for c.Next() {
		//  should just be in the server block once.
		if found {
			return nil, plugin.ErrOnce
		}
		found = true

		// parse the zone list, normalizing each to a FQDN, and
		// using the zones from the server block if none are given.
		args := c.RemainingArgs()
		if len(args) == 0 {
			e.zones = make([]string, len(c.ServerBlockKeys))
			copy(e.zones, c.ServerBlockKeys)
		}
		for _, str := range args {
			e.zones = append(e.zones, plugin.Host(str).Normalize())
		}
	}
	return e, nil

}
