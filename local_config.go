package gopxy

import (
	"encoding/json"
	"io/ioutil"
)

type RemoteConfig struct {
	Host string `json:"host"`
	Code string `json:"code"`
}

func (this *RemoteConfig) NewCopy() *RemoteConfig {
	data := &RemoteConfig{Host: this.Host, Code: this.Code}
	return data
}

type LocalConfig struct {
	RemoteConfigList []RemoteConfig `json:"remote"`
	DefaultCode      string         `json:"default_code"`
	BindHost         string         `json:"bind_host"` //example: 0:0:0:0:8080
}

func Parse(f string) (*LocalConfig, error) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	cfg := &LocalConfig{}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
