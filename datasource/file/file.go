package file

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/stars-palace/statrs-common/pkg/utils/xgo"
	"github.com/stars-palace/statrs-common/pkg/xfile"
	"github.com/stars-palace/statrs-common/pkg/xlogger"
	conf "github.com/stars-palace/statrs-config"
	"github.com/stars-palace/statrs-config/properties"
	"github.com/stars-palace/statrs-config/xini"
	"github.com/stars-palace/statrs-config/yml"
	"log"
	"os"
	"path/filepath"
	"strings"
)

/**
 *
 * Copyright (C) @2020 hugo network Co. Ltd
 * @description
 * @updateRemark
 * @author               hugo
 * @updateUser
 * @createDate           2020/8/20 11:00 上午
 * @updateDate           2020/8/20 11:00 上午
 * @version              1.0
**/
type ConfigType uint32

const (
	//properties 文件
	Properties ConfigType = iota
	//yml 文件
	Yml
	//ini 文件
	InI
)

// fileDataSource file provider.
type fileDataSource struct {
	path        string
	dir         string
	enableWatch bool
	changed     chan struct{}
}

// NewDataSource returns new fileDataSource.
func NewDataSource(path string, watch bool) *fileDataSource {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		logrus.Panicf("new datasource err", err.Error())
	}
	//检查并获取文件夹
	dir := xfile.CheckAndGetParentDir(absolutePath)
	ds := &fileDataSource{path: absolutePath, dir: dir, enableWatch: watch}
	if watch {
		ds.changed = make(chan struct{}, 1)
		xgo.Go(ds.watch)
	}
	return ds
}

// ParseLevel takes a string level and returns the Logrus log level constant.
func ParseConfigType(lvl string) (ConfigType, error) {
	switch strings.ToLower(lvl) {
	case "properties":
		return Properties, nil
	case "yml":
		return Yml, nil
	case "ini":
		return InI, nil
	}

	var l ConfigType
	return l, fmt.Errorf("not a valid config type: %q", lvl)
}

// ReadConfig ...
func (fp *fileDataSource) ReadConfig() ([]byte, error) {
	f, err := os.Open(fp.path)
	//关闭
	defer f.Close()
	//定义返回变量
	if err != nil {
		return nil, err
	}
	fileSuffix := xfile.GetFileSuffix(f)
	contype, err1 := ParseConfigType(fileSuffix)
	if nil != err1 {
		return nil, err1
	}
	switch contype {
	case Properties:
		err = conf.LoadFromFile(f, properties.ReadFileConfig, properties.Unmarshal)
		break
	case Yml:
		err = conf.LoadFromFile(f, yml.ReadFileConfig, yml.Unmarshal)
	case InI:
		err = conf.LoadFromFile(f, xini.ReadFileConfig, xini.Unmarshal)
	default:
		err = errors.New(fmt.Sprintf("can not find config file type %s", fileSuffix))
	}
	if err != nil {
		return nil, err
	}
	return nil, err
}

// Close ...
func (fp *fileDataSource) Close() error {
	close(fp.changed)
	return nil
}

// IsConfigChanged ...
func (fp *fileDataSource) IsConfigChanged() <-chan struct{} {
	return fp.changed
}

// Watch file and automate update.
func (fp *fileDataSource) watch() {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		logrus.Fatal("new file watcher", xlogger.FieldMod("file datasource"), xlogger.Any("err", err.Error()))
	}

	defer w.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-w.Events:
				logrus.Debug("read watch event",
					xlogger.FieldMod("file datasource"),
					xlogger.String("event", filepath.Clean(event.Name)),
					xlogger.String("path", filepath.Clean(fp.path)),
				)
				// we only care about the config file with the following cases:
				// 1 - if the config file was modified or created
				// 2 - if the real path to the config file changed
				const writeOrCreateMask = fsnotify.Write | fsnotify.Create
				if event.Op&writeOrCreateMask != 0 && filepath.Clean(event.Name) == filepath.Clean(fp.path) {
					logrus.Println("modified file: ", event.Name)
					select {
					case fp.changed <- struct{}{}:
					default:
					}
				}
			case err := <-w.Errors:
				// log.Println("error: ", err)
				logrus.Error("read watch error", xlogger.FieldMod("file datasource"), xlogger.Any("err", err.Error()))
			}
		}
	}()

	err = w.Add(fp.dir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
