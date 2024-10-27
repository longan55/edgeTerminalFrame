package core

import (
	"context"
	"regexp"
	"sync"
)

type DataSource interface {
}

type LocalFile struct {
	ConnectMan
	options      *LocalFileOptions
	regexp       *regexp.Regexp
	ctx          context.Context
	cancel       context.CancelFunc
	fileNameChan chan string
	fileDataChan chan []byte
}

func NewLocalFile(options *LocalFileOptions) *LocalFile {
	return &LocalFile{options: options}
}

type LocalFileOptions struct {
	key     string
	name    string
	share   string
	format  string //文件格式：支持 csv tsv 后续支持：xls xlsx
	pattern string //文件名格式
}

func (options *LocalFileOptions) SetKey(key string) {
	options.key = key
}
func (options *LocalFileOptions) GetKey() string {
	return options.key
}

func (options *LocalFileOptions) SetName(share string) {
	options.share = share
}
func (options *LocalFileOptions) GetName() string {
	return options.name
}

func (options *LocalFileOptions) SetShare(key string) {
	options.key = key
}
func (options *LocalFileOptions) GetShare() string {
	return options.share
}

func (options *LocalFileOptions) SetFormat(format string) {
	options.format = format
}
func (options *LocalFileOptions) GetFormat() string {
	return options.format
}

func (options *LocalFileOptions) SetPattern(pattern string) {
	options.pattern = pattern
}
func (options *LocalFileOptions) GetPattern() string {
	return options.pattern
}

// DataSourceContainer 数据源管理容器
type DataSourceContainer struct {
	mux       sync.RWMutex
	container map[string]DataSource
}

func (c *DataSourceContainer) Put(key string, source DataSource) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.container[key] = source
}

// 不存在返回的是nil
func (c *DataSourceContainer) Get(key string) DataSource {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.container[key]
}

func (c *DataSourceContainer) Del(key string) {
	c.mux.Lock()
	defer c.mux.Unlock()
	delete(c.container, key)
}
