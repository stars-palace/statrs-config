package main

/**
 * Copyright (C) @2020 hugo network Co. Ltd
 *
 * @author: hugo
 * @version: 1.0
 * @date: 2020/9/2
 * @time: 19:56
 * @description:
 */

import (
	"fmt"
	conf "github.com/stars-palace/statrs-config"
	"github.com/stars-palace/statrs-config/properties"
	"github.com/stars-palace/statrs-config/xini"
	"github.com/stars-palace/statrs-config/xyml"
	"os"
)

/**
 * Copyright (C) @2020 hugo network Co. Ltd
 *	文件读取测试
 * @author: hugo
 * @version: 1.0
 * @date: 2020/9/2
 * @time: 19:43
 * @description:
 */

type Application struct {
	Name                 string `properties:"brian.application.name" ini:"app.name" yml:"brian.application.name" `                                            //应用名称
	LogLevel             string `properties:"brian.application.log.level" ini:"app.level.log" yml:"brian.application.log.level"`                              // 日志级别
	EnableRpcServer      bool   `properties:"brian.application.enable.RpcServer" ini:"app.RpcServer" yml:"brian.application.enable.RpcServer"`                //是否开启rpc服务
	EnableRegistryCenter bool   `properties:"brian.application.enable.RegistryCenter" ini:"app.RegistryCenter" yml:"brian.application.enable.RegistryCenter"` //是否启用注册中心
	EnableServerClient   bool   `properties:"brian.application.enable.ServerClient" ini:"app.ServerClient" yml:"brian.application.enable.ServerClient"`       //是否开启客户端链接
	RefreshTime          int    `properties:"brian.application.servers.refresh.time" ini:"app.time" yml:"brian.application.servers.refresh.time"`             //本地服务列表刷新时间
}

func main() {
	//TestPropertiesConfig()
	//TestIniConfig()
	TestYmlConfig()
}

// TestYmlConfig 测试yml文件解析
func TestYmlConfig() {
	f, err := os.Open("resources/application.yml")
	defer f.Close()
	if err != nil {
		fmt.Println("read file fail", err)
	}
	err = conf.LoadFromFile(f, xyml.ReadFileConfig, xyml.Unmarshal)
	if err != nil {
		fmt.Println("read file fail", err)
	}
	var app Application
	err = conf.UnmarshalToStruct(&app)
	if err != nil {
		fmt.Println("read file fail", err)
	}
	fmt.Println(app)
}

// TestIniConfig 测试ini文件解析
func TestIniConfig() {
	f, err := os.Open("resources/application.ini")
	defer f.Close()
	if err != nil {
		fmt.Println("read file fail", err)
	}
	err = conf.LoadFromFile(f, xini.ReadFileConfig, xini.Unmarshal)
	if err != nil {
		fmt.Println("read file fail", err)
	}
	var app Application
	err = conf.UnmarshalToStruct(&app)
	if err != nil {
		fmt.Println("read file fail", err)
	}
	fmt.Println(app)
}

// TestPropertiesConfig 测试properties文件解析
func TestPropertiesConfig() {
	f, err := os.Open("resources/application.properties")
	defer f.Close()
	if err != nil {
		fmt.Println("read file fail", err)
	}
	err = conf.LoadFromFile(f, properties.ReadFileConfig, properties.Unmarshal)
	if err != nil {
		fmt.Println("read file fail", err)
	}
	var app Application
	err = conf.UnmarshalToStruct(&app)
	if err != nil {
		fmt.Println("read file fail", err)
	}
	fmt.Println(app)
}
