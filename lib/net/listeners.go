package net

import (
	"bytes"
	"fmt"
	"github.com/drael/GOnetstat"
)

// GetHostFacts gathers facts related to Host
func GetListenersFacts(f Facter) error {
	var facts bytes.Buffer
	t := GOnetstat.Tcp()
	u := GOnetstat.Udp()
	first := true
	for _, p := range t {
		// Check STATE to show only Listening connections
		if p.State == "LISTEN" {
			if first {
				ip_port := fmt.Sprintf("tcp:%v:%v", p.Ip, p.Port)
				facts.WriteString(ip_port)
			} else {
				ip_port := fmt.Sprintf(",tcp:%v:%v", p.Ip, p.Port)
				facts.WriteString(ip_port)
			}
			first = false
		}
	}
	for _, p := range u {
		// Check STATE to show only Listening connections
		if p.State == "CLOSE" {
			if first {
				ip_port := fmt.Sprintf("udp:%v:%v", p.Ip, p.Port)
				facts.WriteString(ip_port)
			} else {
				ip_port := fmt.Sprintf(",udp:%v:%v", p.Ip, p.Port)
				facts.WriteString(ip_port)
			}
			first = false
		}
	}
	f.Add("listeners", facts.String())
	return nil
}
