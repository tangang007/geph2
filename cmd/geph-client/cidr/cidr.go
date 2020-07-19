package cidr

import (
  "net"
  "github.com/yl2chen/cidranger"
)

var CIDRMatcher = cidranger.NewPCTrieRanger()

func AddCIDRIP(ip string) error {
	_, ip_net, error := net.ParseCIDR(ip)
	if error != nil {
		return error
	}
	return AddCIDR(ip_net)
	
}

func AddCIDR(ip_net *net.IPNet) error {
	return CIDRMatcher.Insert(cidranger.NewBasicRangerEntry(*ip_net))
}

func InCIDRList(ip string) (bool, error) {
	return CIDRMatcher.Contains(net.ParseIP(ip))
}
