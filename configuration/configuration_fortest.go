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
				ConfigAMQPUserName: "guest",
				ConfigAMQPPassword: "guest",

				ConfigAMQPHost: "localhost",
				ConfigAMQPPort: "5672",
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
