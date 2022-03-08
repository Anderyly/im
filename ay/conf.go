package ay

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Yaml struct {
	Domain string `yaml:"domain"`
	Mysql  YamlMysql
	Redis  YamlRedis
}

type YamlMysql struct {
	Localhost string `yaml:"localhost"`
	Port      string `yaml:"port"`
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
	Database  string `yaml:"database"`
}

type YamlRedis struct {
	Localhost string `yaml:"localhost"`
	Port      string `yaml:"port"`
	Password  string `yaml:"password"`
}

func (c *Yaml) GetConf() *Yaml {
	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return c
}
