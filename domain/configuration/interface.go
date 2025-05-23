package configuration

type IConfiguration interface {
	GetValue(key string) string
}
