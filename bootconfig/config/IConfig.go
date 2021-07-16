package config

type ConfigValue []byte

func (c *ConfigValue) Bytes() []byte {
	return []byte(*c)
}

func (c *ConfigValue) String() string {
	return string(*c)
}

type IConfig interface {
	Get(string) ([]byte, error)
	GetConfig(string) (ConfigValue, error)
}
