package registry

import (
	"time"
	"github.com/spf13/viper"
)

var v = viper.New()

func Set(key string, value interface{}) {
	v.Set(key, value)
}

func Get(key string) interface{} {
	return v.Get(key)
}

func GetBool(key string) bool {
	return v.GetBool(key)
}

func GetDuration(key string) time.Duration {
	return v.GetDuration(key)
}

func GetFloat64(key string) float64 {
	return v.GetFloat64(key)
}

func GetInt(key string) int {
	return v.GetInt(key)
}

func GetString(key string) string {
	return v.GetString(key)
}

func GetStringMap(key string) map[string]interface{} {
	return v.GetStringMap(key)
}

func GetStringMapString(key string) map[string]string {
	return v.GetStringMapString(key)
}

func GetStringMapStringSlice(key string) map[string][]string {
	return v.GetStringMapStringSlice(key)
}

func GetStringSlice(key string) []string {
	return v.GetStringSlice(key)
}

func GetTime(key string) time.Time {
	return v.GetTime(key)
}
