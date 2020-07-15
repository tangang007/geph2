package extralist

import (
	log "github.com/sirupsen/logrus"
	"github.com/weppos/publicsuffix-go/publicsuffix"
)


func InExtralist(fqdn string) bool {
	dom, err := publicsuffix.Domain(fqdn)
	if err != nil {
		return false
	}
	status := ExtraList[fqdn] || ExtraList[dom]
	if status {
		log.Debugf("%s in extra list", fqdn)
	}
	return ExtraList[fqdn] || ExtraList[dom]
}
