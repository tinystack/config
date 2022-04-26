package config

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cast"
	"github.com/tinystack/errors"
)

type FileType uint8

const (
	UnknownFileType FileType = iota
	YamlFileType
	TomlFileType
	IniFileType
	JsonFileType
	EnvFileType
)

const (
	defaultKeyDelim = "."
)

type Config struct {
	configFile string
	kv         map[string]interface{}
	kvCache    *sync.Map
}

func New() *Config {
	return &Config{
		kv:      make(map[string]interface{}),
		kvCache: new(sync.Map),
	}
}

func (c *Config) GetFloat64(key string) float64 {
	return cast.ToFloat64(c.getValue(key))
}

func (c *Config) GetBool(key string) bool {
	return cast.ToBool(c.getValue(key))
}

func (c *Config) GetString(key string) string {
	return cast.ToString(c.getValue(key))
}

func (c *Config) GetInt(key string) int {
	return cast.ToInt(c.getValue(key))
}

func (c *Config) GetIntSlice(key string) []int {
	return cast.ToIntSlice(c.getValue(key))
}

func (c *Config) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(c.getValue(key))
}

func (c *Config) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(c.getValue(key))
}

func (c *Config) GetStringSlice(key string) []string {
	return cast.ToStringSlice(c.getValue(key))
}

func (c *Config) GetTime(key string) time.Time {
	return cast.ToTime(c.getValue(key))
}

func (c *Config) GetDuration(key string) time.Duration {
	return cast.ToDuration(c.getValue(key))
}

func (c *Config) Get(key string) interface{} {
	return c.getValue(key)
}

func (c *Config) getValue(key string) interface{} {
	lk := strings.ToLower(key)
	if cacheVal, ok := c.kvCache.Load(lk); ok {
		return cacheVal
	}
	keys := strings.Split(lk, defaultKeyDelim)
	val := c.getValueFromMaps(keys)
	c.kvCache.Store(lk, val)
	return val
}

func (c *Config) SetConfigFile(file string) {
	c.configFile = file
}

func (c *Config) LoadConfig() error {
	file, fileType, err := c.getConfigFile()
	if err != nil {
		return err
	}

	b, err := c.readLocalFile(file)
	if err != nil {
		return err
	}

	kv := make(map[string]interface{})
	if err := c.decodeReader(b, kv, fileType); err != nil {
		return err
	}

	mapsKey2Lower(kv)
	c.kv = kv

	return nil
}

func (c *Config) decodeReader(b []byte, cfg map[string]interface{}, fileType FileType) error {
	dc, ok := decoders[fileType]
	if !ok {
		panic(fmt.Sprintf("fileType %v no decoder", fileType))
	}
	return dc.Decode(b, cfg)
}

func (c *Config) getConfigFile() (string, FileType, error) {
	if c.configFile == "" {
		return "", UnknownFileType, errors.New("configFile name is empty")
	}

	index := strings.LastIndex(c.configFile, ".")
	if index == -1 {
		return "", UnknownFileType, errors.New("configFile extension not supported")
	}
	configFileType := c.getFileTypeByExtension(c.configFile[index+1:])
	if configFileType == UnknownFileType {
		return "", UnknownFileType, errors.New("configFile ")
	}

	return c.configFile, configFileType, nil
}

func (c *Config) getFileTypeByExtension(ext string) FileType {
	switch strings.ToLower(ext) {
	case "yaml":
		return YamlFileType
	case "yml":
		return YamlFileType
	case "toml":
		return TomlFileType
	case "ini":
		return IniFileType
	case "json":
		return JsonFileType
	case "env":
		return EnvFileType
	default:
		return UnknownFileType
	}
}

func (c *Config) readLocalFile(file string) ([]byte, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (c *Config) getValueFromMaps(keys []string) interface{} {
	var (
		val  interface{} = c.kv
		nval map[string]interface{}
		ok   bool
	)
	for _, k := range keys {
		nval, ok = val.(map[string]interface{})
		if !ok {
			return nil
		}
		val, ok = nval[k]
		if !ok {
			return nil
		}
	}
	return val
}

func mapsKey2Lower(kv map[string]interface{}) {
	for k, v := range kv {
		switch v.(type) {
		case map[interface{}]interface{}:
			v = cast.ToStringMap(v)
			mapsKey2Lower(v.(map[string]interface{}))
		case map[string]interface{}:
			mapsKey2Lower(v.(map[string]interface{}))
		}
		lk := strings.ToLower(k)
		if lk != k {
			delete(kv, k)
			kv[lk] = v
		}
	}
}
