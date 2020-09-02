package xini

import (
	"errors"
	"fmt"
	"github.com/stars-palace/statrs-common/pkg/xcodec"
	conf "github.com/stars-palace/statrs-config"
	"gopkg.in/ini.v1"
	"os"
	"reflect"
	"strings"
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

const tagName = "ini"

// ReadFileConfig 读取配置
func ReadFileConfig(file *os.File) (map[string]interface{}, error) {
	cfg, err := ini.Load(file) //初始化一个cfg
	if err != nil {
		return nil, err
	}
	//获取所有的selecten
	sections := cfg.Sections()
	return parsingSection(sections)
}

// UnmarshalKey takes a single key and unmarshal it into a Struct.
func Unmarshal(c *conf.Configuration, rawVal interface{}) error {
	//返回对应的类型
	dataType := reflect.TypeOf(rawVal)
	//返回对应类型的reflect.value
	dataValue := reflect.ValueOf(rawVal)
	//判断是否是指针，只有指针才能进行操作
	if dataValue.Kind() == reflect.Ptr {
		//是否时空的
		if dataValue.IsNil() {
			return errors.New("读取配置文件传入的必须是指针")
		}
		// 解引用
		dataValue = dataValue.Elem()
		dataType = dataType.Elem()
	}
	//判断接收数据的类型
	rewType := dataType.Kind()
	switch rewType {
	case reflect.Struct:
		//获取结构体属性的个数
		fieldNum := dataValue.NumField()
		//通过遍历给结构体的属性赋值
		for i := 0; i < fieldNum; i++ {
			field := dataType.Field(i)
			//获取结构体的tag
			tag := field.Tag.Get(tagName)
			//根据名称获取值信息
			fieldValue := dataValue.Field(i)
			if tag == "" {
				c.Mu.RLock()
			}
			//获取配置的值
			value := getValue(c, tag)
			if value == nil {
				continue
			}
			//判断值是否有效。 当值本身非法时，返回 false，例如 reflect Value不包含任何值，值为 nil 等。
			if !fieldValue.IsValid() {
				continue
			}
			if fieldValue.CanInterface() {
				//判断值是否可以被改变
				if fieldValue.CanSet() {
					// TODO 当前只对基本类型处理缺少对结构体中数组和结构体的处理
					switch field.Type.Kind() {
					case reflect.Struct:
						val, err1 := xcodec.UnmarshalStruct(value, field.Type)
						if err1 != nil {
							return err1
						}
						//赋值
						fieldValue.Set(val)
						break
					case reflect.Slice:
						val, err1 := xcodec.UnmarshalArray(value, field.Type)
						if err1 != nil {
							return err1
						}
						//赋值
						fieldValue.Set(val)
						break
					case reflect.Map:
						val, err1 := xcodec.UnmarshalMap(value, dataType.Elem())
						if err1 != nil {
							return err1
						}
						//赋值
						fieldValue.Set(val)
						break
					default:
						//基本本数据类型转换
						val, err1 := xcodec.BasicUnmarshalByType1(value, field.Type)
						if err1 != nil {
							return err1
						}
						//赋值
						fieldValue.Set(reflect.ValueOf(val))
						break
					}
				}

			}
		}
		break
	case reflect.Map:
		resVale, rerr := xcodec.UnmarshalByType(c.KeyMap, dataValue.Elem().Type())
		if nil != rerr {
			return rerr
		}
		//给返回结果赋值
		dataValue.Elem().Set(resVale)
	default:
		return errors.New(fmt.Sprintf("can not Unmarshal config to %s", dataType.String()))
	}
	return nil
}

// 解析值
func getValue(c *conf.Configuration, tag string) interface{} {
	//获取第一个点出现的位置
	i := strings.Index(tag, ".")
	if i >= 0 {
		key := string([]byte(tag)[:i])
		vale := c.Get(key)
		if vale != nil {
			//获取去掉取出来的key
			s := string([]byte(tag)[i+1:])
			return getValueByMap(vale.(map[string]interface{}), s)
		} else {
			return vale
		}
	} else {
		return c.Get(tag)
	}
}

func getValueByMap(data map[string]interface{}, s string) interface{} {
	//获取第一个点出现的位置
	i := strings.Index(s, ".")
	if i >= 0 {
		key := string([]byte(s)[:i])
		vale := data[key]
		if vale != nil {
			//获取去掉取出来的key
			s := string([]byte(s)[i+1:])
			tp := reflect.TypeOf(vale)
			if tp.Kind() == reflect.Map {
				return getValueByMap(vale.(map[string]interface{}), s)
			}
			return nil
		} else {
			return vale
		}
	} else {
		return data[s]
	}
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
			child[v.Name()] = v.Value()
		}
		config[section.Name()] = child
	}
	return config, nil
}
