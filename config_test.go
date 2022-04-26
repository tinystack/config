package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testDataTime time.Time
)

func init() {
	testDataTime, _ = time.ParseInLocation(time.RFC3339, "2022-04-19T13:15:58Z", time.UTC)
}

func testLoadConfig(c *Config, t *testing.T) {
	assert.NoError(t, c.LoadConfig())
	assert.Equal(t, "apps/v1", c.Get("apiVersion"))
	assert.Equal(t, "nginx-test", c.Get("metadata.name"))
	assert.Equal(t, 3.1415926, c.GetFloat64("testdata.pi"))
	assert.Equal(t, true, c.GetBool("testdata.switch"))
	assert.Equal(t, "Deployment", c.GetString("kind"))
	assert.Equal(t, 3, c.GetInt("spec.replicas"))
	assert.Equal(t, []int{2, 8, 16}, c.GetIntSlice("testData.intSlice"))
	assert.Equal(t, map[string]interface{}{
		"name": "nginx-test",
		"lables": map[string]interface{}{
			"app": "nginx",
		},
	}, c.GetStringMap("metadata"))
	assert.Equal(t, map[string]string{
		"app": "nginx",
	}, c.GetStringMapString("spec.selector.matchLabels"))
	assert.Equal(t, []string{"hello", "world"}, c.GetStringSlice("testData.stringSilce"))
	assert.Equal(t, testDataTime, c.GetTime("testData.time"))
	assert.Equal(t, time.Duration(100), c.GetDuration("testData.duration"))
}

func TestYamlLoadConfig(t *testing.T) {
	c := New()
	c.SetConfigFile("./testdata/test.yaml")
	testLoadConfig(c, t)
}

func TestJsonLoadConfig(t *testing.T) {
	c := New()
	c.SetConfigFile("./testdata/test.json")
	testLoadConfig(c, t)
}

func TestTomlLoadConfig(t *testing.T) {
	c := New()
	c.SetConfigFile("./testdata/test.toml")
	testLoadConfig(c, t)
}

func TestIniLoadConfig(t *testing.T) {
	c := New()
	c.SetConfigFile("./testdata/test.ini")
	assert.NoError(t, c.LoadConfig())
	assert.Equal(t, "apps/v1", c.Get("default.apiVersion"))
	assert.Equal(t, 8080, c.GetInt("default.port"))
	assert.True(t, c.GetBool("testData.switch"))
	assert.Equal(t, "Deployment", c.GetString("default.kind"))
	assert.Equal(t, time.Duration(100), c.GetDuration("testData.duration"))
	assert.Equal(t, 3.1415926, c.GetFloat64("testData.pi"))
	assert.Equal(t, testDataTime, c.GetTime("testData.time"))
	assert.Equal(t, map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}, c.GetStringMap("stringMap"))
	assert.Equal(t, map[string]string{
		"key1": "value1",
		"key2": "value2",
	}, c.GetStringMapString("stringMap"))
}

func TestEnvLoadConfig(t *testing.T) {
	c := New()
	c.SetConfigFile("./testdata/test.env")
	assert.NoError(t, c.LoadConfig())
	assert.Equal(t, "apps/v1", c.Get("apiversion"))
	assert.Equal(t, testDataTime, c.GetTime("time"))
	assert.True(t, c.GetBool("switch"))
	assert.Equal(t, "Deployment", c.GetString("kind"))
	assert.Equal(t, 3.1415926, c.GetFloat64("pi"))
	assert.Equal(t, time.Duration(100), c.GetDuration("duration"))
	assert.Equal(t, 8080, c.GetInt("port"))
}
