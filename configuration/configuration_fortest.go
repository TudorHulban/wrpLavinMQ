package configuration

type ConfigurationTest struct {
	configuration map[string]string
}

var _ IConfiguration = ConfigurationTest{}

// error added for same signature with other configurations.
func NewConfigurationTest() (*ConfigurationTest, error) {
	return &ConfigurationTest{
			configuration: map[string]string{
				ConfigAMQPProtocol: "amqp",
				ConfigAMQPNameUser: "gtest",
				ConfigAMQPPassword: "gtest",

				ConfigAMQPHost: "localhost",
				ConfigAMQPPort: "5672",

				ConfigAMQPVirtualHost: "gtest_host",

				ConfiqAMQPNameExchange:        "ex12345",
				ConfiqAMQPNameQueueMessages:   "q12345",
				ConfiqAMQPNameQueueAggregates: "q67890",
			},
		},
		nil
}

func (config ConfigurationTest) GetConfigurationValue(key string) string {
	result, exists := config.configuration[key]
	if exists {
		return result
	}

	return ""
}
