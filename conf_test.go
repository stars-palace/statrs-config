package conf

import (
	"github.com/stars-palace/statrs-config/properties"
	"os"
	"testing"
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
	Name                 string `properties:"brian.application.name"`                  //应用名称
	LogLevel             string `properties:"brian.application.log.level"`             // 日志级别
	EnableRpcServer      bool   `properties:"brian.application.enable.RpcServer"`      //是否开启rpc服务
	EnableRegistryCenter bool   `properties:"brian.application.enable.RegistryCenter"` //是否启用注册中心
	EnableServerClient   bool   `properties:"brian.application.enable.ServerClient"`   //是否开启客户端链接
	RefreshTime          int    `properties:"brian.application.servers.refresh.time"`  //本地服务列表刷新时间
}

func TestPropertiesConfig(t *testing.T) {
	f, err := os.Open("../resources/application.properties")
	defer f.Close()
	if err != nil {
		t.Error("read file fail", err)
	}
	err = LoadFromFile(f, properties.ReadFileConfig, properties.Unmarshal)
	if err != nil {
		t.Error("read file fail", err)
	}
	var app Application
	err = UnmarshalToStruct(&app)
	if err != nil {
		t.Error("read file fail", err)
	}
	t.Skip(app)
}
