# mindbenders

Logger

``` go
type ILogWriter interface {
    WriteLogs(context.Context, logrus.Fields, logrus.Level, string)
}

//ILogger ...
type IDotpeLogger interface {
    ILogWriter
    GinLogger() gin.HandlerFunc
}
```

In order to initialize one has to have pass the 2 things for sure.

1. `logger.WithAppInfo(os.Getenv("APP"), os.Getenv("CWD"), os.Getenv("ENV"))`, definition is as follows

    ```go
    func WithAppInfo(app, env, wd string) Option {
        return func(dlogger *dlogger) {
            dlogger.app = app
            dlogger.env = env
            dlogger.wd = wd
        }
    }
    ```

2. `logger.WithHookContainer(hookContainer)`, definition is as follows

    ```go
    //.....
    type IHookContainer interface {
        GetHook() (logrus.Hook, error)
    }

    //.....
    func WithHookContainer(hookContainer IHookContainer) Option {
        hook, err := hookContainer.getHook()
        if err != nil {
            return nil
        }
        return WithHook(hook)
    }
    ```

Hook

1. Elastic hook container

    ```go
    logconf = logger.NewKibanaConfig(url, key, secret, os.Getenv("APP"), "")
    ```

2. File Hook container

    ```go
    logconf = logger.NewKibanaConfig(url, key, secret, os.Getenv("APP"), "")
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
var hookContainer logger.IHookContainer

// InitLogger ..
DLogger, err := logger.Init(
    logger.WithAppInfo(os.Getenv("APP"), os.Getenv("CWD"), os.Getenv("ENV")),
    logger.WithHookContainer(hookContainer),
    logger.WithAccessLogOptions(
        logger.AccessLogOptionRequestBody,
        aopt1, // check definition of aopt1, similary you can pass more functions as you need
    ),
    logger.WithLogOptions(opt1),// check definition of opt1, similary you can pass more functions as you need
)

func aopt1(c *gin.Context, fields *logrus.Fields) {
    c.Set("ip", c.ClientIP())
}

func opt1(ctx context.Context, fields *logrus.Fields) {
    (*fields)["clientIP"] = ctx.Value("ip")
}
```

attaching GinLogger to add accessLog
`apiGroup.Use(utils.DLogger.GinLogger())`

Recovery Middleware

The call `ginmiddleware.Recovery(utils.DLogger)` is for making the gin-engine failure safe.
> **Note:** you can't control crashes in orphened go-routines
