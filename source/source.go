package source

import (
	"github.com/pkg/errors"
	"github.com/varyoo/albumpact/meta"
	"net/url"
	"regexp"
	"strings"
)

func init() {
	sources = make(map[string]*source)
}

var sources map[string]*source

type getter func(url string) (meta.Album, []meta.Track, error)
type source struct {
	get getter
}

func RegisterSource(host string, get getter) {
	sources[host] = &source{get}
}
func hosts() []string {
	hosts := make([]string, 0)
	for h, _ := range sources {
		hosts = append(hosts, h)
	}
	return hosts
}
func Get(rawURL string) (meta.Album, []meta.Track, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, nil, errors.Wrap(err, "url parsing")
	}
	r := regexp.MustCompile(`([^\.]+\.[^\.]*)$`)
	domain := r.FindStringSubmatch(u.Host)
	if len(domain) != 2 {
		return nil, nil, errors.New("not a valid URL")
	}
	host := domain[1]
	s := sources[host]
	if s == nil {
		return nil, nil, errors.Errorf("unsupported host '%s', available: '%s'",
			host, strings.Join(hosts(), ","))
	}
	a, l, err := s.get(rawURL)
	if err != nil {
		err = errors.Wrap(err, host)
	}
	return a, l, err
}
