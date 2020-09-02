package yml

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

/**
 * Copyright (C) @2020 hugo network Co. Ltd
 * 解析yml文件依赖yml2
 * @author: hugo
 * @version: 1.0
 * @date: 2020/9/1
 * @time: 23:19
 * @description:
 */

func ReadFileConfig(file *os.File) (map[string]interface{}, error) {
	config := make(map[string]interface{})
	fd, err := ioutil.ReadAll(file)
	if nil != err {
		return nil, err
	}
	if err := yaml.Unmarshal(fd, &config); err != nil {
		return nil, err
	}
	return config, nil
}
