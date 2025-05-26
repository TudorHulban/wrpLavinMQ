package configuration

type IConfiguration interface {
	GetConfigurationValue(key string) string
}
