package config

type Config struct {
	ConfigDataBase struct {
		Host         string `envconfig:"HOST_DB"`
		Port         string `envconfig:"PORT_DB"`
		User         string `envconfig:"USER_DB"`
		Password     string `envconfig:"PASSWORD_DB"`
		NameDataBase string `envconfig:"NAME_DB"`
	}
	ListenPort int `envconfig:"LISTEN_PORT"`
}
