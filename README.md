# config
Go config Package

## 示例

### 加载配置文件
```go
import "github.com/tinystack/config"

cfg := config.New()
cfg.SetConfigFile("./config.yaml")

if err := cfg.LoadConfig(); err != nil {
	panic("config load failed")
}
```

### 可用的API
```go
cfg.GetFloat64()
cfg.GetBool()
cfg.GetString()
cfg.GetInt()
cfg.GetIntSlice()
cfg.GetStringMap()
cfg.GetStringMapString()
cfg.GetStringSlice()
cfg.GetTime()
cfg.GetDuration()
cfg.Get()
```

## 支持的配置文件类型

- ini
- yaml/yml
- toml
- json
- env