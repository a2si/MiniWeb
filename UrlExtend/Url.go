package UrlExtend

import (
	"net/url"
	"strings"

	DevLogs "github.com/a2si/MiniWeb/DevLogs"
)

func (self *TUrl) SetUrl(URL string) {
	DevLogs.Debug("TUrl.SetUrl")
	var err error
	self.URL, err = url.Parse(URL)
	if err != nil {
		DevLogs.Warn("TUrl.SetUrl Err=" + err.Error())
	}
}

func (self *TUrl) AddParam(Name string, Value string) {

}

func (self *TUrl) GenSign(Name string) {

}

func (self *TUrl) SetUserPassword(Name string, Password string) {
	self.URL.User = url.UserPassword(Name, Password)
}

func (self *TUrl) SetScheme(Scheme string) {
	self.URL.Scheme = Scheme
}

func (self *TUrl) GetScheme() string {
	return self.URL.Scheme
}

func (self *TUrl) IsTls() bool {
	if strings.ToLower(self.URL.Scheme) == "https" {
		return true
	}
	if self.URL.Port() == "443" {
		return true
	}
	return false
}

func (self *TUrl) GetPath() string {
	return self.URL.Path
}

func (self *TUrl) GetHost() string {
	return self.URL.Hostname()
}

func (self *TUrl) GetPort() string {
	var (
		Port string = self.URL.Port()
	)
	if len(Port) == 0 {
		if self.IsTls() == true {
			Port = "443"
		} else {
			Port = "80"
		}
	}
	return Port
}

/*
	Scheme     string
	Opaque     string    // encoded opaque data
	Host       string    // host or host:port
	Path       string    // path (relative paths may omit leading slash)
	RawPath    string    // encoded path hint (see EscapedPath method)
	ForceQuery bool      // append a query ('?') even if RawQuery is empty
	RawQuery   string    // encoded query values, without '?'
	Fragment   string    // fragment for references, without '#'
*/

func (self *TUrl) GetEncode() string {
	var (
		urlValue url.Values = urlParseQuery(self.URL.RawQuery)
		Query    string     = urlValue.Encode()
	)
	return Query
}

func urlParseQuery(query string) url.Values {
	m := make(url.Values)
	parseQuery(m, query)
	return m
}

func parseQuery(m url.Values, query string) (err error) {
	for query != "" {
		key := query
		if i := strings.IndexAny(key, "&"); i >= 0 {
			key, query = key[:i], key[i+1:]
		} else {
			query = ""
		}
		if key == "" {
			continue
		}
		value := ""
		if i := strings.Index(key, "="); i >= 0 {
			key, value = key[:i], key[i+1:]
		}
		key, err1 := url.QueryUnescape(key)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		value, err1 = url.QueryUnescape(value)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		m[key] = append(m[key], value)
	}
	return err
}
