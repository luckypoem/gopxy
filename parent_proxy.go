package gopxy

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net"
	"net/http"
)

type ParentProxy struct {

}

func(this *ParentProxy) processConn(conn net.Conn) {
	defer conn.Close()
	br := bufio.NewReader(conn)
	request, err := http.ReadRequest(br)
	if err != nil {
		logrus.Errorf("Read request fail, err:%v", err)
		return
	}
	url := fmt.Sprintf("%s://%s", request.URL.Scheme, request.URL.Hostname())
	cli := http.Client{}
	proxyUrl := "https://orange-cake-2c2f.xxxtest.workers.dev"
	dataBuf := make([]byte, 4096)
	buf := bytes.NewBuffer(dataBuf)
	request.Write(buf)
	io.MultiReader(buf, br)
	req, _ := http.NewRequest("POST", proxyUrl, br)
	req.Header.Add("code", "haha")
	req.Header.Add("m_proxy_target", url)
	req.Header.Add("m_proxy_method", request.Method)
	rsp, err := cli.Do(req)
	logrus.Tracef("Do request to target:%s, err:%v", proxyUrl, err)
	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		logrus.Errorf("Read data from proxy fail, err:%v", err)
		return
	}
	logrus.Infof("data:%s", string(data))
}

func(this *ParentProxy) Start() {
	acc, err := net.Listen("tcp", "127.0.0.1:1087")
	if err != nil {
		panic(err)
	}
	for {
		cli, err := acc.Accept()
		if err != nil {
			logrus.Errorf("Accept fail, err:%+v", err)
			continue
		} else {
			go this.processConn(cli)
		}
	}
}
