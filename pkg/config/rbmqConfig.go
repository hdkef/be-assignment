package config

type RBMQConfig struct {
	RBMQURL string
}

func InitRBMQConfig() *RBMQConfig {
	config := &RBMQConfig{
		RBMQURL: getEnvOrPanic("RABBITMQ_URL"),
	}
	return config
}
