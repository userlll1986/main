package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// 定义Config结构体
type Config struct {
	// 服务器相关配置
	Server struct {
		Part   string `yml:"part"`
		Model  string `yml:"model"`
		Banner struct {
			Name    string `yml:"name"`
			Loading bool   `yml:"loading"`
		} `yml:"banner"`
	} `yml:"server"`
	// 数据库相关配置
	Data struct {
		Category string `yml:"category"`
		Prefix   string `yml:"prefix"`
		DataBase string `yml:"database"`
		Ip       string `yml:"ip"`
		Part     string `yml:"part"`
		Username string `yml:"username"`
		Password string `yml:"password"`
		Sql      bool   `yml:"sql"`
		Init     struct {
			Name   string `yml:"name"`
			Status bool   `yml:"status"`
		} `yml:"init"`
		Test struct {
			Name   string `yml:"name"`
			Status bool   `yml:"status"`
		} `yml:"test"`
	} `yml:"data"`
}

func NewConfig() *Config {
	return &Config{}
}

// 读取配置映射到结构体
func (config *Config) ReadConfig() *Config {

	file, err := ioutil.ReadFile("config/config.yml")
	if err != nil {
		log.Fatalln("读取文件config.yml发生错误")
		return nil
	}
	if yaml.Unmarshal(file, config) != nil {
		log.Fatalln("解析文件config.yml发生错误")
		return nil
	}
	return config
}
