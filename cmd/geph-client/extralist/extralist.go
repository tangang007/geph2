package extralist

import (
	"io/ioutil"
	"bufio"
	"os"
	"regexp"
	"strings"
	"net/http"
	log "github.com/sirupsen/logrus"
	"github.com/miekg/dns"
	"github.com/geph-official/geph2/cmd/geph-client/cidr"
)

var SourceConfigs map[string]ListSource

var ExtraList map[string]bool


func ExtralistFilter(content []byte, pattern regexp.Regexp) (string, error) {
	log.Info("IN ExtralistFilter")
	var result strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {
		line := scanner.Text()
		//log.Info("LINE:", line)
		//log.Info("PATTERN:", pattern)
		matches := pattern.FindStringSubmatch(line)
		//log.Info("MATCHES", matches)
		if len(matches) != 0 {
			result.WriteString(matches[1])
			result.WriteString("\n")
			//log.Info("%v got now matched patterns!", line)
		} else {
			log.Infof("%v got no matched pattern!", line)
		}
	}
	return result.String(), nil
}

func LoadExtralistContent(content []byte) error {
	lines := strings.Split(string(content), "\n")
	ExtraList = make(map[string]bool, len(lines))
	log.Infof("Extralist contains %v domain(s)", len(lines))
	for _, domain := range lines {
		if len(strings.TrimSpace(domain)) == 0 {
			continue // ignore empty lines
		}

		if _, ok := dns.IsDomainName(domain); !ok {
			log.Infof("%v is not a valid domain name, try CIDR mode", domain)
			err := cidr.AddCIDRIP(domain)
			if err  != nil {
				log.Infof("%v is not a valid CIDR, ignored", domain)
			} else {
				log.Debugf("Adding [%v] to locallist", domain)
			}
		} else {
			if _, ok := ExtraList[domain]; ok {
				log.Debugf("%v already exists in cache", domain)
			}
			ExtraList[domain] = true
		}
		
	}
	return nil
}

func LoadExtralist(source string) error {	
	content, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}
	return LoadExtralistContent(content)
}

func UpdateExtraList(url string, target string, pattern *regexp.Regexp,client *http.Client) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var parsed string

	if pattern != nil{
		parsed, err = ExtralistFilter(content, *pattern)
	}

	
	log.Infof("Parsed %v", url)
	if err != nil {
		return err
	}

	if target != "" {
		dst, err := os.OpenFile(target, os.O_TRUNC|os.O_WRONLY, 0666)
		if os.IsNotExist(err) {
			dst, err = os.Create(target)
			if err != nil {
				log.Error("Failed to create ", target)
				return err
			}
		}
		
		if err != nil {
			return err
		}
		_, err = dst.Write([]byte(parsed))
		dst.Close()
	} else {
		LoadExtralistContent([]byte(parsed))
	}
	if err != nil {
		return err
	}
	return err
}
