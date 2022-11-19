package bootconfig

import (
	"log"

	"gitlab.com/dotpe/mindbenders/bootconfig/confmanager"
)

var ConfManager confmanager.IConfigExt

func Init(env string) error {
	var conf1 confmanager.IConfig
	var err error
	if env == "dev" {
		conf1, err = confmanager.GetFileConfigManager()
		log.Println(err)
	} else {
		conf1 = confmanager.GetSecretManager(env)
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
	confmanager.IConfig
}

func (cfgmgr *conf) MustGet(key string) []byte {
	raw, err := cfgmgr.IConfig.Get(key)
	if err != nil {
		log.Fatalf("unable to read config: %v\b", err)
	}
	return raw
}

func (cfgmgr *conf) MustGetGlobal(key string) []byte {
	raw, err := cfgmgr.IConfig.GetGlobal(key)
	if err != nil {
		log.Fatalf("unable to read config: %v\b", err)
	}
	return raw
}
