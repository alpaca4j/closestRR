package closestRR

//https://tools.ietf.org/html/rfc1035
//RFC8482 - no more ANY
import (
	"context"
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/plugin/pkg/nonwriter"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"net"
)

var log = clog.NewWithPlugin("closestRR")

type closestRR struct {
	Next  plugin.Handler
	zones []string
}

func (e closestRR) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	log.Debug("Zone is: ", state.Name())

	// If the zone does not match one of ours, just pass it on.
	if plugin.Zones(e.zones).Matches(state.Name()) == "" {
		return plugin.NextOrFailure(e.Name(), e.Next, ctx, w, r)
	}

	// If the query is not for A type record, just pass it on.
	log.Debug("The question is: ", state.QType())
	if state.QType() != dns.TypeA {
		return plugin.NextOrFailure(e.Name(), e.Next, ctx, w, r)
	}

	nw := nonwriter.New(w)

	rcode, err := plugin.NextOrFailure(e.Name(), e.Next, ctx, nw, r)
	if err != nil {
		// Simply return if there was an error.
		return rcode, err
	}

	w.WriteMsg(e.filterOutOneRecordBySubnet(nw.Msg, state.IP()))
	return rcode, err

}

func (e *closestRR) filterOutOneRecordBySubnet(m *dns.Msg, ip string) *dns.Msg {
	cidr := "/23"
	_, network, _ := net.ParseCIDR(ip + cidr)
	log.Debug("Requestor IP is: " + ip)
	log.Debug("Calculated network with ", cidr, " equals: ", network)

	if len(m.Answer) <= 1 {
		return m
	}

	var newAnswer []dns.RR
	for _, a := range m.Answer {

		if a.Header().Rrtype == dns.TypeA {
			recordA := a.(*dns.A).A.String()
			isIpInNetwork := network.Contains(net.ParseIP(recordA))
			log.Debug("Does this network: ", network, " contain this ip ", net.ParseIP(recordA), " ? The anwser is: ", isIpInNetwork)

			if isIpInNetwork {
				newAnswer = append(newAnswer, a)
				log.Debug("Found a matching address, returning: ", a.String())
			}
		} else {
			newAnswer = append(newAnswer, a)
		}
	}
	if len(newAnswer) > 0 {
		m.Answer = newAnswer
		return m
	}

	return m
}

// Name implements the Handler interface.
func (e closestRR) Name() string { return "closestRR" }
