# mindbenders

Logger

``` go
type ILogWriter interface {
	WriteLogs(context.Context, Fields, Level, string)
	Info(context.Context, Fields, string)
	Error(context.Context, Fields, string)
	Warn(context.Context, Fields, string)
	Debug(context.Context, Fields, string)
}

//ILogger ...
type IDotpeLogger interface {
    ILogWriter
    GinLogger() gin.HandlerFunc
}
```

In order to initialize one has to have pass the 2 things for sure.

1. `logger.WithAppInfo(os.Getenv("APP"))`, definition is as follows

    ```go
    func WithAppInfo(app string) Option {
        return func(dlogger *dlogger) {
            dlogger.app = app
        }
    }
    ```

**logger initializations**
Usage [[link](https://gitlab.com/dotcomino/2c/-/blob/master/utils/logger.go)]

```go
package utils

import (
    "os"

    mbinterfaces "gitlab.com/dotpe/mindbenders/interfaces"
    logger "gitlab.com/dotpe/mindbenders/logging"
)

var DLogger mbinterfaces.IDotpeLogger

//write your logic to initialize or fetch hookContainer objectct

// InitLogger ..
DLogger, err := logger.Init(
    logger.WithAppInfo(os.Getenv("APP")),
    logger.WithAccessLogOptions(
        logger.AccessLogOptionRequestBody,
        aopt1, // check definition of aopt1, similary you can pass more functions as you need
    ),
    logger.WithLogOptions(opt1),// check definition of opt1, similary you can pass more functions as you need
)

func aopt1(c *gin.Context, fields *Fields) {
    c.Set("ip", c.ClientIP())
}

func opt1(ctx context.Context, fields *Fields) {
    (*fields)["clientIP"] = ctx.Value("ip")
}
```

attaching GinLogger to add accessLog
`apiGroup.Use(utils.DLogger.GinLogger())`

Recovery Middleware

The call `ginmiddleware.Recovery(utils.DLogger)` is for making the gin-engine failure safe.
> **Note:** you can't control crashes in orphened go-routines
