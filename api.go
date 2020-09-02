package conf

/**
 *
 * Copyright (C) @2020 hugo network Co. Ltd
 * @description
 * @updateRemark
 * @author               hugo
 * @updateUser
 * @createDate           2020/9/1 11:00 上午
 * @updateDate           2020/9/1 11:00 上午
 * @version              1.0
**/

import (
	"io"
	"os"
)

/**
 * Copyright (C) @2020 hugo network Co. Ltd
 *
 * @author: hugo
 * @version: 1.0
 * @date: 2020/8/2
 * @time: 12:32
 * @description:
 */

// DataSource ...配置数据的来源
type DataSource interface {
	ReadConfig() ([]byte, error)
	IsConfigChanged() <-chan struct{}
	io.Closer
}

//默认配置
var defaultConfiguration = New()

// Unmarshaller ...
type Unmarshaller = func([]byte, interface{}) error

// UnmarshallerToStruct  ... 读取配置文件到结构体中
type UnmarshallerToStruct = func(*Configuration, interface{}) error

// 读取配置值文件
type ReadConfigFile = func(file *os.File) (map[string]interface{}, error)

// 从数据原中读取配置
// LoadFromDataSource load configuration from data source
// if data source supports dynamic config, a monitor goroutinue
// would be
func LoadFromDataSource(ds DataSource, unmarshaller Unmarshaller) error {
	return defaultConfiguration.LoadFromDataSource(ds, unmarshaller)
}

// Load loads configuration from provided provider with default defaultConfiguration.
func LoadFromReader(r io.Reader, unmarshaller Unmarshaller) error {
	return defaultConfiguration.LoadFromReader(r, unmarshaller)
}

// LoadFromFile 从文件中读取配置文件
func LoadFromFile(file *os.File, readConfig ReadConfigFile, unmarshaller UnmarshallerToStruct) error {
	return defaultConfiguration.LoadFromFile(file, readConfig, unmarshaller)
}

// Apply ...
func Apply(conf map[string]interface{}) error {
	return defaultConfiguration.apply(conf)
}

// Get returns an interface. For a specific value use one of the Get____ methods.
func Get(key string) interface{} {
	return defaultConfiguration.Get(key)
}

// Set set config value for key
func Set(key string, val interface{}) {
	defaultConfiguration.Set(key, val)
}
