package bootconfig

type iConfig interface {
	// To retreive *env* specific config
	// `prod/redis` or `stage1/redis`
	Get(string) ([]byte, error)

	// To retreive Global config | non env specific
	// `dynamoDB`
	GetGlobal(string) ([]byte, error)
}

type ConfigManager interface {
	iConfig
	//This panics if no such config found
	MustGet(string) []byte

	//This panics if no such global config found
	MustGetGlobal(string) []byte
}
