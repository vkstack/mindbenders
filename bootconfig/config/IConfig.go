package config

type IConfig interface {
	// To retreive *env* specific config
	// `prod/redis` or `stage1/redis`
	Get(string) ([]byte, error)

	// To retreive Global config
	// `dynamoDB`
	GetGlobal(string) ([]byte, error)
}
