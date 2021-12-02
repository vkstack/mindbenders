package config

type ConfigValue []byte

func (c *ConfigValue) Bytes() []byte {
	return []byte(*c)
}

func (c *ConfigValue) String() string {
	return string(*c)
}

type IConfig interface {
	// env specific config
	Get(string) ([]byte, error)

	// env specific config
	GetConfig(string) (ConfigValue, error)

	// env specifig config
	GetGlobal(string) ([]byte, error)
}
