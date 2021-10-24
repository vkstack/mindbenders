package logging

import (
	"os"
	"sync"

	"gitlab.com/dotpe/mindbenders/interfaces"
)

var lock sync.Mutex
var logger interfaces.IDotpeLogger
var host, _ = os.Hostname()

//InitLogger sets up the logger object with LoeggerOptions provided.
//It returns reference logger object and error
func Init(opts ...Option) (interfaces.IDotpeLogger, error) {
	if logger == nil {
		lock.Lock()
		defer lock.Unlock()
		if logger == nil {
			var dlogger = new(dlogger)
			for _, opt := range opts {
				opt(dlogger)
			}
			err := dlogger.finalizeEssentials()
			if err != nil {
				return nil, err
			}
			logger = dlogger
		}
	}
	return logger, nil
}
