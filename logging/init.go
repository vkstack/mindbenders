package logging

import (
	"io/ioutil"
	"sync"

	"github.com/sirupsen/logrus"
	"gitlab.com/dotpe/mindbenders/interfaces"
)

var lock = &sync.Mutex{}
var logger interfaces.IDotpeLogger

//InitLogger sets up the logger object with LoeggerOptions provided.
//It returns reference logger object and error
func InitLogger(lops *LoggerOptions) (interfaces.IDotpeLogger, error) {
	if logger == nil {
		lock.Lock()
		if logger == nil {
			if err := initlogger(lops); err != nil {
				return nil, err
			}
		}
		lock.Unlock()
	}
	return logger, nil
}

func initlogger(lops *LoggerOptions) error {
	log := logrus.New()
	log.SetNoLock()
	hook, err := lops.IConfig.getHook()
	if err != nil {
		log.Panic(err)
		return err
	}
	log.Hooks.Add(hook)
	if lops.LOGENV != "dev" {
		log.Out = ioutil.Discard
	}
	logger = &dlogger{Logger: log, Lops: *lops}
	return nil
}
