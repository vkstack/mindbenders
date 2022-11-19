package confmanager

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"gitlab.com/dotpe/mindbenders/errors"
)

const devconfig = "config.local.json"

// FileConfig It is for dev configs only.
// For staging & prod SecretManager will be used.
type fileConfig struct {
	keyVal map[string]interface{}
}

func GetFileConfigManager() (IConfig, error) {
	cfgMgr := fileConfig{}
	localConfig, err := os.Open(path.Join(os.Getenv("CONFIGDIR"), devconfig))
	if err != nil {
		return nil, errors.WrapMessage(err, "config file not opening")
	}
	rawbytes, err := ioutil.ReadAll(localConfig)
	if err != nil {
		return nil, errors.WrapMessage(err, "couldn't read from rawbytes")
	}
	err = json.Unmarshal(rawbytes, &cfgMgr.keyVal)
	if err != nil {
		return nil, errors.WrapMessage(err, "couldn't unmarshall rawbytes")
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
	return nil, errors.New("no such config exists:\t" + key)
}

func (cfgmgr *fileConfig) GetConfig(key string) (raw []byte, err error) {
	if raw, err = cfgmgr.Get(key); err != nil {
		return nil, err
	}
	return raw, nil
}

// Global config
func (cfgmgr *fileConfig) GetGlobal(key string) (raw []byte, err error) {
	return cfgmgr.Get(key)
}
