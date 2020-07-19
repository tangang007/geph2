package locallist

import (
	log "github.com/sirupsen/logrus"
	"github.com/geph-official/geph2/cmd/geph-client/cidr"
)

var LocalIPCIDR = []string{
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
	"fd00::/8",
}

var is_initialed = false

func InitCIDRWithLocal() {
	if !is_initialed {
		for _,ip := range LocalIPCIDR {
			log.Debugf("Adding [%v] to locallist", ip)
			cidr.AddCIDRIP(ip)
		}
		is_initialed = true
	}
}

func IsLocalList(target string) (bool, error) {
	if ! is_initialed {
		InitCIDRWithLocal()
	}
	return cidr.InCIDRList(target)
}
