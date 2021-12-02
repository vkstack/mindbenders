package bootconfig

import (
	"log"

	"gitlab.com/dotpe/mindbenders/bootconfig/config"
	"gitlab.com/dotpe/mindbenders/bootconfig/confmanager"
)

var ConfManager config.IConfig

func Init(env string) error {
	var err error
	if env == "dev" {
		ConfManager, err = confmanager.GetFileConfigManager()
		log.Println(err)
	} else {
		ConfManager = confmanager.GetSecretManager(env)
	}
	return err
}
