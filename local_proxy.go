package gopxy

import (
	"bytes"
	"github.com/elazarl/goproxy"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net/http"
)

func init() {
	logrus.SetLevel(logrus.TraceLevel)
}

type LocalProxy struct {
	cfg *LocalConfig
}


func New(cfg *LocalConfig) *LocalProxy {
	return &LocalProxy{cfg}
}

func(this *LocalProxy) Start() error {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		url := req.URL.Host

		if len(req.URL.Path) != 0 {
			url = url + req.URL.Path
		}
		if len(req.URL.RawQuery) != 0 {
			url = url + "?" + req.URL.RawQuery
		}
		cli := http.Client{}
		proxyUrl := "https://orange-cake-2c2f.xxxtest.workers.dev/" + url
		var dataBuf []byte
		buf := bytes.NewBuffer(dataBuf)
		req.Write(buf)
		mt := io.MultiReader(buf, req.Body)
		mreq, _ := http.NewRequest(req.Method, proxyUrl, mt)
		rsp, err := cli.Do(mreq)
		logrus.Tracef("Do request to proxy:%s, err:%v, url:%s, path:%s, query:%s", proxyUrl, err, url, req.URL.Path, req.URL.RawQuery)
		return mreq, rsp
	})
	log.Fatal(http.ListenAndServe(":8080", proxy))
	return nil
}