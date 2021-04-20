package bootconf

import (
	"gitlab.com/dotpe/mindbenders/bootconf/config"
	"gitlab.com/dotpe/mindbenders/bootconf/confmanager"
)

var ConfManager config.IConfig

func Init(env string) error {
	var err error
	if env == "dev" {
		ConfManager, err = confmanager.GetFileConfigManager()
	} else {
		ConfManager = confmanager.GetSecretManager(env)
	}
	return err
}
