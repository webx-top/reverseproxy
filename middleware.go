package reverseproxy

import (
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/admpub/log"
	"github.com/webx-top/echo"
	. "github.com/webx-top/reverseproxy/log"
)

var ErrNoBackends = errors.New("no backends")
var ErrAllBackendsDead = errors.New("all backends are dead")

type ProxyOptions struct {
	PathPrefix      string   //网址路径前缀，符合这个前缀的将反向代理到其它服务器
	Engine          string   //反向代理使用的引擎，值为fast时使用fastHTTP，否则使用标准HTTP
	Hosts           []string //支持通过反向代理访问的后台服务器集群，例如 192.168.0.2:8080
	FlushInterval   time.Duration
	DialTimeout     time.Duration
	RequestTimeout  time.Duration
	RequestIDHeader string
	ResponseBefore  func(Context) bool
	ResponseAfter   func(Context) bool
	router          *ProxyRouter
}

func (p *ProxyOptions) AddHost(hosts ...string) {
	if p.router == nil {
		p.router = NewProxyRouter(p.Hosts...)
	}
	if len(hosts) > 0 {
		p.router.AddHost(hosts...)
	}
}

func Proxy(options *ProxyOptions) echo.MiddlewareFunc {
	options.AddHost()
	config := &ReverseProxyConfig{
		FlushInterval:        options.FlushInterval,
		DialTimeout:          options.DialTimeout,
		RequestTimeout:       options.RequestTimeout,
		RequestIDHeader:      options.RequestIDHeader,
		ResponseBefore:       options.ResponseBefore,
		ResponseAfter:        options.ResponseAfter,
		DisabledAloneService: true,
		Router:               options.router,
	}
	config.DisabledAloneService = true
	var rPxy ReverseProxy
	if strings.ToLower(options.Engine) == `fast` {
		rPxy = &FastReverseProxy{PassingBrowsingURL: true}
		options.Engine = `FastHTTP`
	} else {
		rPxy = &NativeReverseProxy{PassingBrowsingURL: true}
		options.Engine = `Standard`
	}
	_, err := rPxy.Initialize(*config)
	if err != nil {
		panic(err.Error())
	}
	prefixLength := len(options.PathPrefix)
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			urlPath := c.Request().URL().Path()
			if len(urlPath) > prefixLength && urlPath[0:prefixLength] == options.PathPrefix {
				/*
					if options.router.hostNum < 1 {
						return ErrAllBackendsDead
					}
				*/
				rPxy.HandlerForEcho(c.Response(), c.Request())
				return nil
			}
			return h.Handle(c)
		})
	}
}

func NewProxyRouter(hosts ...string) *ProxyRouter {
	pr := &ProxyRouter{
		onlineHosts:  []string{},
		offlineHosts: []string{},
	}
	if len(hosts) > 0 {
		pr.AddHost(hosts...)
	}
	return pr
}

type ProxyRouter struct {
	dst                string //目标网址
	resultHost         string //最终操作的主机
	resultReqData      *RequestData
	resultIsDead       bool
	logEntry           *LogEntry
	onlineHosts        []string
	offlineHosts       []string
	hostNum            int
	DisabledLogRequest bool
}

func (r *ProxyRouter) AddHost(hosts ...string) *ProxyRouter {
	for _, host := range hosts {
		for idx, offlineHost := range r.offlineHosts {
			if offlineHost == host {
				r.offlineHosts = append(r.offlineHosts[0:idx], r.offlineHosts[idx+1:]...)
			}
		}
	}
	r.onlineHosts = append(r.onlineHosts, hosts...)
	r.hostNum = len(r.onlineHosts)
	return r
}

func (r *ProxyRouter) ChooseBackend(host string) (rd *RequestData, err error) {
	if r.hostNum < 1 {
		err = ErrNoBackends
		r.onlineHosts = r.offlineHosts
		r.offlineHosts = []string{}
		r.hostNum = len(r.onlineHosts)
	}
	r.resultHost = host
	rd = &RequestData{
		Backend:    r.dst,
		BackendIdx: 0,
		BackendKey: host,
		BackendLen: r.hostNum,
		Host:       host,
		StartTime:  time.Now(),
	}
	idx := 0
	if r.hostNum > 1 {
		idx = rand.Intn(r.hostNum - 1)
	}
	rd.BackendIdx = idx
	rd.BackendKey = r.onlineHosts[idx]
	rd.Backend = `http://` + r.onlineHosts[idx]
	return
}

func (r *ProxyRouter) EndRequest(reqData *RequestData, isDead bool, fn func() *LogEntry) error {
	r.resultReqData = reqData
	r.logEntry = fn()
	r.resultIsDead = isDead
	if r.resultIsDead {
		r.offlineHosts = append(r.offlineHosts, r.onlineHosts[r.resultReqData.BackendIdx])
		r.onlineHosts = append(r.onlineHosts[0:r.resultReqData.BackendIdx], r.onlineHosts[r.resultReqData.BackendIdx+1:]...)
		r.hostNum = len(r.onlineHosts)
	}

	if !r.DisabledLogRequest {
		log.Infof("== Request: %7s %s => Completed %d in %vs", r.logEntry.Method, r.logEntry.Path, r.logEntry.StatusCode, r.logEntry.TotalDuration.Seconds())
	}
	return nil
}
