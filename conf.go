package conf

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/stars-palace/statrs-common/pkg/xcast"
	"github.com/stars-palace/statrs-common/pkg/xfile"
	"github.com/stars-palace/statrs-common/pkg/xmap"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"sync"
)

/**
 * Copyright (C) @2020 hugo network Co. Ltd
 *
 * @author: hugo
 * @version: 1.0
 * @date: 2020/8/2
 * @time: 12:27
 * @description:
 */

// Configuration provides configuration for application.
var configType = "properties"

// 解析器，给结构体赋值
var unmarshallerToStruct UnmarshallerToStruct

//配置整个系统的应用
type Configuration struct {
	Mu       sync.RWMutex
	override map[string]interface{}
	keyDelim string

	KeyMap    *sync.Map
	onChanges []func(*Configuration)

	watchers map[string][]func(*Configuration)
}

//获取配置文件类型
func GetConfigType() string {
	return configType
}

//设置配置文件类型
func SetConfigType(v string) {
	configType = v
}

//放在默认的配置中
// UnmarshalKey takes a single key and unmarshal it into a Struct with default defaultConfiguration.
func UnmarshalKey(key string, rawVal interface{}, opts ...GetOption) error {
	//配置默认设置
	return defaultConfiguration.UnmarshalKey(key, rawVal, opts...)
}

//解析配置
func Unmarshal(rawVal interface{}) error {
	//配置默认设置
	return nil //defaultConfiguration.Unmarshal(rawVal)
}

// UnmarshalToStruct 解析配置到结构体
func UnmarshalToStruct(rawVal interface{}) error {
	//配置默认设置
	return defaultConfiguration.localUnmarshalToStruct(rawVal)
}

//解析配置
func (c *Configuration) localUnmarshalToStruct(rawVal interface{}) error {
	return unmarshallerToStruct(c, rawVal)
}

const (
	defaultKeyDelim = "."
)

// ErrInvalidKey ...
var ErrInvalidKey = errors.New("invalid key, maybe not exist in config")

// New constructs a new Configuration with provider.
func New() *Configuration {
	return &Configuration{
		override:  make(map[string]interface{}),
		keyDelim:  defaultKeyDelim,
		KeyMap:    &sync.Map{},
		onChanges: make([]func(*Configuration), 0),
		watchers:  make(map[string][]func(*Configuration)),
	}
}

// LoadFromDataSource ...
func (c *Configuration) LoadFromDataSource(ds DataSource, unmarshaller Unmarshaller) error {
	content, err := ds.ReadConfig()
	if err != nil {
		return err
	}

	if err := c.Load(content, unmarshaller); err != nil {
		return err
	}

	go func() {
		for range ds.IsConfigChanged() {
			if content, err := ds.ReadConfig(); err == nil {
				_ = c.Load(content, unmarshaller)
				for _, change := range c.onChanges {
					change(c)
				}
			}
		}
	}()

	return nil
}

// Load ...
func (c *Configuration) Load(content []byte, unmarshaller Unmarshaller) error {
	configuration := make(map[string]interface{})
	if err := unmarshaller(content, &configuration); err != nil {
		return err
	}
	return c.apply(configuration)
}

// Load loads configuration from provided data source.
func (c *Configuration) LoadFromReader(reader io.Reader, unmarshaller Unmarshaller) error {
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	return c.Load(content, unmarshaller)
}

// Load loads configuration from config file.
func (c *Configuration) LoadFromFile(file *os.File, readConfig ReadConfigFile, unmarshaller UnmarshallerToStruct) error {
	configuration, err := readConfig(file)
	if err != nil {
		return err
	}
	unmarshallerToStruct = unmarshaller
	configType = xfile.GetFileSuffix(file)
	return c.apply(configuration)
}

// Set ...
func (c *Configuration) Set(key string, val interface{}) error {
	paths := strings.Split(key, c.keyDelim)
	lastKey := paths[len(paths)-1]
	m := deepSearch(c.override, paths[:len(paths)-1])
	m[lastKey] = val
	return c.apply(m)
	// c.KeyMap.Store(key, val)
}

//应用配置
func (c *Configuration) apply(conf map[string]interface{}) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	var changes = make(map[string]interface{})

	xmap.MergeStringMap(c.override, conf)
	for k, v := range c.traverse(c.keyDelim) {
		orig, ok := c.KeyMap.Load(k)
		if ok && !reflect.DeepEqual(orig, v) {
			changes[k] = v
		}
		c.KeyMap.Store(k, v)
	}

	if len(changes) > 0 {
		c.notifyChanges(changes)
	}

	return nil
}

// UnmarshalKey takes a single key and unmarshal it into a Struct.
func (c *Configuration) UnmarshalKey(key string, rawVal interface{}, opts ...GetOption) error {
	var options = defaultGetOptions
	for _, opt := range opts {
		opt(&options)
	}

	config := mapstructure.DecoderConfig{
		DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
		Result:     rawVal,
		TagName:    options.TagName,
	}
	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return err
	}
	if key == "" {
		c.Mu.RLock()
		defer c.Mu.RUnlock()
		return decoder.Decode(c.override)
	}

	value := c.Get(key)
	if value == nil {
		return errors.Wrap(ErrInvalidKey, key)
	}

	return decoder.Decode(value)
}

//通知改变的进行改变
func (c *Configuration) notifyChanges(changes map[string]interface{}) {
	var changedWatchPrefixMap = map[string]struct{}{}

	for watchPrefix := range c.watchers {
		for key := range changes {
			// 前缀匹配即可
			// todo 可能产生错误匹配
			if strings.HasPrefix(key, watchPrefix) {
				changedWatchPrefixMap[watchPrefix] = struct{}{}
			}
		}
	}

	for changedWatchPrefix := range changedWatchPrefixMap {
		for _, handle := range c.watchers[changedWatchPrefix] {
			go handle(c)
		}
	}
}

//深度搜索
func deepSearch(m map[string]interface{}, path []string) map[string]interface{} {
	for _, k := range path {
		m2, ok := m[k]
		if !ok {
			m3 := make(map[string]interface{})
			m[k] = m3
			m = m3
			continue
		}
		m3, ok := m2.(map[string]interface{})
		if !ok {
			m3 = make(map[string]interface{})
			m[k] = m3
		}
		m = m3
	}
	return m
}

// Get returns the value associated with the key
func (c *Configuration) Get(key string) interface{} {
	return c.find(key)
}

//查找
func (c *Configuration) find(key string) interface{} {
	//直接读取
	dd, ok := c.KeyMap.Load(key)
	if ok {
		return dd
	}

	paths := strings.Split(key, c.keyDelim)
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	m := xmap.DeepSearchInMap(c.override, paths[:len(paths)-1]...)
	dd = m[paths[len(paths)-1]]
	c.KeyMap.Store(key, dd)
	return dd
}

//遍历
func (c *Configuration) traverse(sep string) map[string]interface{} {
	data := make(map[string]interface{})
	lookup("", c.override, data, sep)
	return data
}

func lookup(prefix string, target map[string]interface{}, data map[string]interface{}, sep string) {
	for k, v := range target {
		pp := fmt.Sprintf("%s%s%s", prefix, sep, k)
		if prefix == "" {
			pp = k
		}
		if dd, err := xcast.ToStringMapE(v); err == nil {
			lookup(pp, dd, data, sep)
		} else {
			data[pp] = v
		}
	}
}

func ChickHavePoint(s string) bool {
	if len(s) > 0 {
		i := strings.Index(s, ".")
		if i >= 0 {
			return true
		}
	}
	return false
}
