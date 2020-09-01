package yml

import (
	"fmt"
	"io/ioutil"
	"os"
)

/**
 * Copyright (C) @2020 hugo network Co. Ltd
 * 一直是map
 * @author: hugo
 * @version: 1.0
 * @date: 2020/9/1
 * @time: 23:19
 * @description:
 */

func Read() {
	f, err := os.Open("../resources/application.yml")
	if err != nil {
		fmt.Println("read file fail", err)
	}
	defer f.Close()

	fd, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("read to fd fail", err)
	}

	fmt.Println(string(fd))
}
