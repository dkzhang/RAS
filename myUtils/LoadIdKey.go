package myUtils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func LoadIdKey(filename string) (idKey *IdKeyConfig, err error) {
	idKey = &IdKeyConfig{}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read IdKeyConfig file error: %v", err)
	}

	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, idKey)
	if err != nil {
		return nil, fmt.Errorf("unmarshal IdKeyConfig file error: %v", err)
	}

	return idKey, nil
}

type IdKeyConfig struct {
	SecretId  string
	SecretKey string
}
