package extralist

// Parsing extralist config file

import (
	_ "io/ioutil"
	"time"
	"os"
	log "github.com/sirupsen/logrus"
	"regexp"
	 "github.com/go-ini/ini"
)

func ParseConfigFile(path string) error {
	cfg, err := ini.LooseLoad(path)
	if err != nil {
		return err
	}

	sources := cfg.Sections()
	
	for _, source := range sources {
		log.Infoln("Reading source ", source.Name())
		res, err := ParseSource(source)
		if err != nil {
			return err
		}
		if SourceConfigs == nil {
			SourceConfigs = make(map[string]ListSource)
		}
		SourceConfigs[source.Name()] = res
	}

	return nil
}

func ParseSource(source *ini.Section) (ListSource,error) {
	var src ListSource
	if source.HasKey("url") && source.Key("url") != nil {
		log.Infof("Source [%v] will be updated from:[%v]", source.Name(), source.Key("url"))
		src.url = source.Key("url").String()

		// regex pattern only affects if url exists
		if source.HasKey("pattern") && source.Key("pattern") != nil {
			log.Infof("Source [%v] will be formatted by regexp: [%v]", source.Name(), source.Key("pattern"))
			src.pattern = regexp.MustCompile(source.Key("pattern").Value())
		}
	} else {
		src.url = ""
	}

	if source.HasKey("mode") {
		if source.Key("mode").String() == "CIDR" {
			src.cidr = true
		} else {
			src.cidr = false
		}
	}

	if source.HasKey("interval") {
		log.Infof("Source [%v] will be cached to [%v]", source.Name(), source.Key("file"))
		interval, err := source.Key("interval").Int()
		if err != nil {
			log.Error("Error while parsing interval in ", source.Name())
			return src, err
		}
		src.interval = time.Duration(interval) * time.Hour
	} else {
		src.interval = 24 * time.Hour // default interval is 24h
	}

	log.Infof("Will updating [%v] every %v hours", source.Name(), src.interval.Hours())

	if source.HasKey("file") {
		log.Infof("Source [%v] will be cached to [%v]", source.Name(), source.Key("file"))
		src.dst = source.Key("file").String()
		_, err := os.Stat(src.dst)
		if err != nil {
			if os.IsNotExist(err) {
				log.Infof("File [%v] does not exists, creating a new one", src.dst)
				_, err := os.Create(src.dst)
				if err != nil {
					log.Error("Failed to create file:", src.dst)
					return src, err
				}
			} else {
				return src, err
			}
		}
	}
	
	return src,nil
}
