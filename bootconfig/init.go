package bootconfig

import (
	"log"
)

var ConfManager ConfigManager

func Init(env string) error {
	var conf1 iConfig
	var err error
	if env == "dev" {
		conf1, err = GetFileConfigManager()
		log.Println(err)
	} else {
		conf1 = GetSecretManager(env)
	}
	ConfManager = &conf{conf1}
	return err
}

func MustInit(env string) {
	if err := Init(env); err != nil {
		log.Fatalf("unable to initialize config manager: %v\b", err)
	}
}

type conf struct {
	iConfig
}

func (cfgmgr *conf) MustGet(key string) []byte {
	raw, err := cfgmgr.iConfig.Get(key)
	if err != nil {
		log.Fatalf("unable to read config: %v\b", err)
	}
	return raw
}

func (cfgmgr *conf) MustGetGlobal(key string) []byte {
	raw, err := cfgmgr.iConfig.GetGlobal(key)
	if err != nil {
		log.Fatalf("unable to read config: %v\b", err)
	}
	return raw
}
