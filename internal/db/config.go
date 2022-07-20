package db

import (
	_ "embed"
	"gopkg.in/yaml.v3"
	"log"
)

//go:embed config.yaml
var taskYaml string

type TaskExternalConfig struct {
	Dsn0           string                   `yaml:"dsn0"`
	Dsn1           string                   `yaml:"dsn1"`
	T0             string                   `yaml:"t0"`
	T1             string                   `yaml:"t1"`
	DsnType0       string                   `yaml:"dsnType0"`
	DsnType1       string                   `yaml:"dsnType1"`
	StartCondition map[string][]interface{} `yaml:"startCondition"`
	StartFuncType  int                      `yaml:"startFuncType"`
	NextFuncType   int                      `yaml:"nextFuncType"`
	NextKey        string                   `yaml:"nextKey"`
	EndFuncType    int                      `yaml:"endFuncType"`
	EndKey         string                   `yaml:"endKey"`
	EndVal         interface{}              `yaml:"endVal"`
	DeleteKey      string                   `yaml:"deleteKey"`
	Task           *Task                    `yaml:"task"`
	Created        bool                     `yaml:"created"`
}

func GetConfig() (taskExternalConfigList []*TaskExternalConfig, err error) {
	m := make(map[interface{}]interface{})
	if err = yaml.Unmarshal([]byte(taskYaml), &m); err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println(m["task"])

	var taskYamlBytes []byte
	if taskYamlBytes, err = yaml.Marshal(m["task"]); err != nil {
		return nil, err
	}

	log.Println(string(taskYamlBytes))

	if err = yaml.Unmarshal(taskYamlBytes, &taskExternalConfigList); err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(taskExternalConfigList[0].Dsn0)
	return
}
