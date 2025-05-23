package configurationtest

import (
	"github.com/TudorHulban/wrpLavinMQ/domain/configuration"
)

type ConfigurationTest struct {
	configuration map[string]string
}

var _ configuration.IConfiguration = ConfigurationTest{}

func NewConfigurationTest() (*ConfigurationTest, error) {
	return &ConfigurationTest{
			configuration: map[string]string{
				configuration.ConfigAMQPProtocol:  "amqp",
				configuration.ConfigAMQPUserName:  "guest",
				configuration.ConfigAMQPPasswword: "guest",

				configuration.ConfigAMQPHost: "localhost",
				configuration.ConfigAMQPPort: "5672",
			},
		},
		nil
}

func (config ConfigurationTest) GetValue(key string) string {
	result, exists := config.configuration[key]
	if exists {
		return result
	}

	return ""
}
