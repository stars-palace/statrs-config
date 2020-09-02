package xyml

import (
	"fmt"
	"os"
	"testing"
)

/**
 * Copyright (C) @2020 hugo network Co. Ltd
 *
 * @author: hugo
 * @version: 1.0
 * @date: 2020/9/1
 * @time: 23:22
 * @description:
 */

func TestReadFileConfig(t *testing.T) {
	f, err := os.Open("../resources/application.yml")
	defer f.Close()
	if err != nil {
		fmt.Println("read file fail", err)
	}
	config, err1 := ReadFileConfig(f)
	if err1 != nil {
		fmt.Println("read config fail", err)
	}
	fmt.Println(config)
}
