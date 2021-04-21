package config

import (
	"Ranking-of-Internet-Corp-By-Go/entity"
	"github.com/go-yaml/yaml"
	"io/ioutil"
)

var Cfg *entity.Config = nil

/**
解析配置文件
 */
func ParseCfgFile() error {
	if Cfg == nil {
		cfgFile, err := ioutil.ReadFile("./config/config.yaml")
		if err != nil {
			return err
		}
		Cfg = &entity.Config{}
		if err := yaml.Unmarshal(cfgFile, Cfg); err != nil {
			return err
		}
	}
	return nil
}