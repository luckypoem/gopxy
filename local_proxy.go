package gopxy

import (
	"bytes"
	"fmt"
	"github.com/elazarl/goproxy"
	"github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

type LocalProxy struct {
	cfg    *LocalConfig
	logger *logrus.Logger
	sess   int32
}

func New(cfg *LocalConfig, logger *logrus.Logger) *LocalProxy {
	if cfg == nil {
		panic("nil config")
	}
	return &LocalProxy{cfg: cfg, logger: logger, sess: 0}
}

func (this *LocalProxy) randomHost() *RemoteConfig {
	if len(this.cfg.RemoteConfigList) == 0 {
		panic("no config found..")
	}
	host := this.cfg.RemoteConfigList[rand.Int()%len(this.cfg.RemoteConfigList)].NewCopy()
	if len(host.Code) == 0 && len(this.cfg.DefaultCode) != 0 {
		host.Code = this.cfg.DefaultCode
	}
	return host
}

func (this *LocalProxy) proxyfunc(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	rnHost := this.randomHost()
	sessionid := atomic.AddInt32(&this.sess, 1)
	this.logger.Infof("Recv request to host:%s, schema:%s, path:%s, query:%s, use config:%+v to proxy it. sessionid:%d",
		req.URL.Hostname(), req.URL.Scheme, req.URL.Path, req.URL.RawQuery, *rnHost, sessionid)
	urlsuffix := req.URL.Host
	if len(req.URL.Path) != 0 {
		urlsuffix = urlsuffix + req.URL.Path
	}
	if len(req.URL.RawQuery) != 0 {
		urlsuffix = urlsuffix + "?" + req.URL.RawQuery
	}
	cli := http.Client{}
	proxyUrl := fmt.Sprintf("https://%s/%s", rnHost.Host, urlsuffix)
	var dataBuf []byte
	buf := bytes.NewBuffer(dataBuf)
	_ = req.Write(buf)
	mt := io.MultiReader(buf, req.Body)
	mreq, err := http.NewRequest(req.Method, proxyUrl, mt)
	if err != nil {
		this.logger.Errorf("Create proxy req fail, err:%v, session:%d, proxy url:%s", err, sessionid, proxyUrl)
	}
	for k, v := range req.Header {
		mreq.Header.Add(k, strings.Join(v, "; "))
	}
	this.logger.Tracef("Proxy req info:%+v, sessionid:%d", *req, sessionid)
	mreq.Header.Add("__m_proxy_schema", req.URL.Scheme)
	mreq.Header.Add("__m_proxy_host", req.URL.Hostname())
	mreq.Header.Add("__m_proxy_referer", req.Referer())
	mreq.Header.Add("__m_proxy_check_code", rnHost.Code)
	left := time.Now()
	rsp, err := cli.Do(mreq)
	cost := time.Now().Sub(left)
	if err != nil {
		this.logger.Errorf("Do request to proxy:%+v net fail, err:%v, url:%s, sessionid:%d, cost:%d", *rnHost, err, urlsuffix, sessionid, cost/time.Millisecond)
	} else {
		this.logger.Infof("Do request to proxy:%+v finish, url:%s, sessionid:%d, status:%d, cost:%d", *rnHost, urlsuffix, sessionid, rsp.StatusCode, cost/time.Millisecond)
	}
	return mreq, rsp
}

func (this *LocalProxy) Start() error {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest().DoFunc(this.proxyfunc)
	err := http.ListenAndServe(this.cfg.BindHost, proxy)
	return err
}
