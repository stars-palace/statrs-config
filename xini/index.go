package xini

import (
	"errors"
	"gopkg.in/ini.v1"
	"os"
)

/**
 *
 * Copyright (C) @2020 hugo network Co. Ltd
 * @description
 * @updateRemark
 * @author               hugo
 * @updateUser
 * @createDate           2020/9/2 1:46 下午
 * @updateDate           2020/9/2 1:46 下午
 * @version              1.0
**/

func ReadFileConfig(file *os.File) (map[string]interface{}, error) {
	cfg, err := ini.Load(file) //初始化一个cfg
	if err != nil {
		return nil, err
	}
	//获取所有的selecten
	sections := cfg.Sections()
	return parsingSection(sections)
}

// parsingSection 解释ini的 Sections形成map
func parsingSection(sections []*ini.Section) (map[string]interface{}, error) {
	config := make(map[string]interface{})
	if len(sections) <= 0 {
		return nil, errors.New("sections Length cannot be zero ")
	}
	for _, section := range sections {
		child := make(map[string]interface{})
		keys := section.Keys()
		for _, v := range keys {
			child[v.String()] = v.Value()
		}
		config[section.Name()] = child
	}
	return config, nil
}
