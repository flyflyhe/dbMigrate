package db

import (
	_ "embed"
	"gopkg.in/yaml.v3"
	"log"
)

//go:embed config.yaml
var TaskYaml string

type TaskExternalConfig struct {
	Name           string                   `json:"name"`
	Dsn0           string                   `json:"dsn0" yaml:"dsn0"`
	Dsn1           string                   `json:"dsn1" yaml:"dsn1"`
	T0             string                   `json:"t0" yaml:"t0"`
	T1             string                   `json:"t1" yaml:"t1"`
	DsnType0       string                   `json:"dsnType0" yaml:"dsnType0"`
	DsnType1       string                   `json:"dsnType1" yaml:"dsnType1"`
	StartCondition map[string][]interface{} `json:"startCondition" yaml:"startCondition"`
	StartFuncType  int                      `json:"startFuncType" yaml:"startFuncType"`
	NextFuncType   int                      `json:"nextFuncType" yaml:"nextFuncType"`
	NextKey        string                   `json:"nextKey" yaml:"nextKey"`
	EndFuncType    int                      `json:"endFuncType" yaml:"endFuncType"`
	EndKey         string                   `json:"endKey" yaml:"endKey"`
	EndVal         interface{}              `json:"endVal" yaml:"endVal"`
	DeleteKey      string                   `json:"deleteKey" yaml:"deleteKey"`
	Created        bool                     `json:"created" yaml:"created"`
}

func GetConfig(yamlBytes []byte) (taskExternalConfigList []*TaskExternalConfig, err error) {
	m := make(map[interface{}]interface{})
	if err = yaml.Unmarshal(yamlBytes, &m); err != nil {
		log.Println(err)
		return nil, err
	}

	var taskYamlBytes []byte
	if taskYamlBytes, err = yaml.Marshal(m["task"]); err != nil {
		return nil, err
	}

	log.Println(string(taskYamlBytes))

	if err = yaml.Unmarshal(taskYamlBytes, &taskExternalConfigList); err != nil {
		log.Println(err)
		return nil, err
	}
	return
}
