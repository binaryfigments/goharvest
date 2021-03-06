package emailmx

import (
	"github.com/miekg/dns"
	"golang.org/x/net/idna"
	"golang.org/x/net/publicsuffix"
)

// Data struct
type Data struct {
	Domain       string     `json:"domain,omitempty"`
	Records      []*Records `json:"records,omitempty"`
	Error        string     `json:"error,omitempty"`
	ErrorMessage string     `json:"errormessage,omitempty"`
}

type Records struct {
	Server     string `json:"server,omitempty"`
	Preference uint16 `json:"preference,omitempty"`
}

func Get(domain string, nameserver string) *Data {
	r := new(Data)

	domain, err := idna.ToASCII(domain)
	if err != nil {
		r.Error = "Failed"
		r.ErrorMessage = err.Error()
		return r
	}

	// Validate
	domain, err = publicsuffix.EffectiveTLDPlusOne(domain)
	if err != nil {
		r.Error = "Failed"
		r.ErrorMessage = err.Error()
		return r
	}

	r.Domain = domain
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeMX)
	m.MsgHdr.RecursionDesired = true
	c := new(dns.Client)
	in, _, err := c.Exchange(m, nameserver+":53")
	if err != nil {
		r.Error = "Failed"
		r.ErrorMessage = err.Error()
		return r
	}
	for _, ain := range in.Answer {
		if a, ok := ain.(*dns.MX); ok {
			records := new(Records)
			records.Server = a.Mx
			records.Preference = a.Preference
			r.Records = append(r.Records, records)
		}
	}
	return r
}
