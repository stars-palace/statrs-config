package properties

import (
	"fmt"
	"os"
	"testing"
)

/**
 * properties文件读取
 * Copyright (C) @2020 hugo network Co. Ltd
 * @description
 * @updateRemark
 * @author               hugo
 * @updateUser
 * @createDate           2020/9/2 1:36 下午
 * @updateDate           2020/9/2 1:36 下午
 * @version              1.0
**/
func TestReadFileConfig(t *testing.T) {
	f, err := os.Open("../resources/application.properties")
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
