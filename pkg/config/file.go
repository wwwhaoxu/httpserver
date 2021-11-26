package config

import (
	"gopkg.in/yaml.v2"
	"io"
	"k8s.io/klog/v2"
	"os"
)

type confing struct {
	path string
	conf map[interface{}]interface{}
}

func NewConfig(conf map[interface{}]interface{}, path string) *confing {
	return &confing{
		conf: conf,
		path: path,
	}
}

func (c *confing) LoadFile() error {
	file, err := os.Open(c.path)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &c.conf)
	if err != nil {
		klog.Errorf("Unmarshal: %v", err)
	}
	return nil
}
