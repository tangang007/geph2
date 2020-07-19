package extralist

import (
	"time"
	"regexp"
)

type list_interface interface {
	Url() string
	Pattern() *regexp.Regexp
	Dst() string
	Bypass() bool
	Interval() time.Duration
} 

type ListSource struct {
	url string // update upstream url
	pattern regexp.Regexp // patterns for everyline in upstream contents
	dst string // destination file path
	bypass bool // bypass or proxy
	interval time.Duration
}


func (src *ListSource) Url() string {
	return src.url
}


func (src *ListSource) Pattern() *regexp.Regexp {
	return &src.pattern
}


func (src *ListSource) Dst() string {
	return src.dst
}


func (src *ListSource) Bypass() bool {
	return src.bypass
}


func (src *ListSource) Interval() *time.Duration {
	return &src.interval
}
