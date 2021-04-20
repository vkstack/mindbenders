package confmanager

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"

	"gitlab.com/dotpe/mindbenders/bootconfig/config"
)

const devconfig = "config.local.json"

//FileConfig It is for dev configs only.
//For staging & prod SecretManager will be used.
type fileConfig struct {
	keyVal map[string]interface{}
}

func GetFileConfigManager() (config.IConfig, error) {
	cfgMgr := fileConfig{}
	localConfig, err := os.Open(path.Join(os.Getenv("cwd"), os.Getenv("CONFIGDIR"), devconfig))
	if err != nil {
		return nil, err
	}
	rawbytes, err := ioutil.ReadAll(localConfig)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rawbytes, &cfgMgr.keyVal)
	if err != nil {
		return nil, err
	}
	//can have some preprocessing logic
	return &cfgMgr, nil
}

func (cfgmgr *fileConfig) Get(key string) ([]byte, error) {
	if val, ok := cfgmgr.keyVal[key]; ok {
		byteVal, err := json.Marshal(val)
		if err != nil {
			return nil, errors.New("error in Unmarshal config")
		}
		return byteVal, nil
	}
	//handle this error very carefully
	return nil, errors.New("no such config exists")
}
