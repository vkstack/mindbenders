package logging

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/olivere/elastic/v7/aws"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	awsauth "github.com/smartystreets/go-aws-auth"
	"github.com/snowzach/rotatefilehook"
	"gopkg.in/sohlich/elogrus.v7"
)

const (
	// LOGPANIC level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	LOGPANIC logrus.Level = iota
	// LOGFATAL level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	LOGFATAL
	// LOGERROR level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	LOGERROR
	// LOGWARN level. Non-critical entries that deserve eyes.
	LOGWARN
	// LOGINFO level. General operational entries about what's going on inside the
	// application.
	LOGINFO
	// LOGDEBUG level. Usually only enabled when debugging. Very verbose logging.
	LOGDEBUG
	// LOGTRACE level. Designates finer-grained informational events than the Debug.
	LOGTRACE

	timeFormat = "02/Jan/2006:15:04:05 -0700"
)

//DPLogger ...
type DPLogger struct {
	Lops   LoggerOptions
	Logger *logrus.Logger
}

//KibanaConfig Mandatory for kibana logging
type KibanaConfig struct {
	Client,
	AccessKey,
	SecretKey,
	APPID,
	Hostname string
}

//LoggerOptions is set of config data for logg
type LoggerOptions struct {
	KibanaConfig
	APP, // Service
	APPID, // Service application ID
	LOGENV, // Dev/Debug/Production
	WD string // Working directory of the application
	COREL interface{}
}

//WriteLogs ...
func (dLogger *DPLogger) WriteLogs(ctx context.Context, fields logrus.Fields, cb logrus.Level, MessageKey string, args ...interface{}) {
	if ctx == nil {
		return
	}

	pc, file, line, _ := runtime.Caller(1)
	_, funcname := filepath.Split(runtime.FuncForPC(pc).Name())
	file = strings.ReplaceAll(file, dLogger.Lops.WD, "")
	corRelationID := ctx.Value(dLogger.Lops.COREL).(map[string]interface{})
	for idx := range fields {
		switch fields[idx].(type) {
		case int8, int16, int32, int64, int,
			uint8, uint16, uint32, uint64, uint,
			float32, float64,
			string, bool:
		default:
			tmp, _ := json.Marshal(fields[idx])
			fields[idx] = string(tmp)
		}
	}
	for idx := range args {
		if idx == 5 {
			break
		}
		fields[fmt.Sprintf("field_%d", idx)] = args[idx]
	}
	fields["caller"] = fmt.Sprintf("%s:%d\n%s", file, line, funcname)
	fields["appid"] = dLogger.Lops.APPID
	fields["requestID"] = corRelationID["requestID"]
	fields["sessionID"] = corRelationID["sessionID"]
	entry := dLogger.Logger.WithFields(fields)
	switch cb {
	case LOGERROR:
		entry.Error(MessageKey)
	case LOGWARN:
		entry.Warn(MessageKey)
	case LOGINFO:
		entry.Info(MessageKey)
	case LOGFATAL:
		entry.Log(logrus.FatalLevel, MessageKey)
	case LOGPANIC:
		entry.Log(logrus.PanicLevel, MessageKey)
	default:
		entry.Fatal(MessageKey)
	}
}

//InitLogger ...
func InitLogger(lops *LoggerOptions) (*DPLogger, error) {
	if lops.Hostname == "" {
		if x, err := os.Hostname(); err != nil {
			lops.Hostname = "unknown"
		} else {
			lops.Hostname = x
		}
	}
	if lops == nil {
		return nil, errors.New("invalid logger options")
	}
	var hook logrus.Hook
	var err error
	log := logrus.New()
	log.SetNoLock()
	if lops.LOGENV == "DEV" {
		formatter := &logrus.TextFormatter{
			ForceColors:               false,
			DisableColors:             false,
			EnvironmentOverrideColors: false,
			DisableTimestamp:          false,
			FullTimestamp:             false,
			TimestampFormat:           "",
			DisableSorting:            false,
			SortingFunc:               nil,
			DisableLevelTruncation:    false,
			QuoteEmptyFields:          false,
			FieldMap:                  nil,
			CallerPrettyfier:          nil,
		}
		hook, err = rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
			Filename:   "logfile.log",
			MaxSize:    5,
			MaxBackups: 7,
			MaxAge:     7,
			Level:      logrus.DebugLevel,
			Formatter:  formatter,
		})
		if err != nil {
			fmt.Println(err)
			log.Panic(err)
			return nil, err
		}
	} else {
		client, err := newElasticClient(&lops.KibanaConfig)
		if err != nil {
			fmt.Println(err)
			log.Panic(err)
			return nil, err
		}
		hook, err = elogrus.NewAsyncElasticHook(client, "", logrus.DebugLevel, lops.APP)
		if err != nil {
			log.Panic(err)
			return nil, err
		}
	}
	log.Hooks.Add(hook)
	if lops.LOGENV != "DEV" {
		log.Out = ioutil.Discard
	}
	return &DPLogger{Logger: log, Lops: *lops}, nil
}

func newElasticClient(kibops *KibanaConfig) (*elastic.Client, error) {
	if kibops.Client == "" {
		log.Fatal("missing -client-url KIBANA")
	}
	if kibops.AccessKey == "" {
		log.Fatal("missing -access-key or AWS_ACCESS_KEY environment variable")
	}
	if kibops.SecretKey == "" {
		log.Fatal("missing -secret-key or AWS_SECRET_KEY environment variable")
	}

	sniff := flag.Bool("sniff", false, "Enable or disable sniffing")

	flag.Parse()
	log.SetFlags(0)

	signingClient := aws.NewV4SigningClient(awsauth.Credentials{
		AccessKeyID:     kibops.AccessKey,
		SecretAccessKey: kibops.SecretKey,
	})

	client, err := elastic.NewClient(
		elastic.SetURL(kibops.Client),
		elastic.SetSniff(*sniff),
		elastic.SetHealthcheck(*sniff),
		elastic.SetHttpClient(signingClient),
	)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	fmt.Println("AWS ElasticSearchConnection succeeded")
	return client, nil
}

//GinLogger ...
func (dLogger *DPLogger) GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// other handler can change c.Path so:
		path := c.Request.URL.Path
		start := time.Now()
		var sessionID, requestID string
		if _, ok := c.Request.Header["Session_id"]; !ok || len(c.Request.Header["Session_id"]) == 0 { // Handling OPTIONS request
			sessionID = "unknownSession"
		} else {
			sessionID = c.Request.Header["Session_id"][0]
		}
		if _, ok := c.Request.Header["Request_id"]; !ok || len(c.Request.Header["Request_id"]) == 0 { // Handling OPTIONS request
			requestID = xid.New().String()
		} else {
			requestID = c.Request.Header["Request_id"][0]
		}
		c.Set("requestID", requestID)
		c.Set("sessionID", sessionID)
		ctx := context.WithValue(c, dLogger.Lops.COREL,
			map[string]interface{}{
				"requestID": requestID,
				"sessionID": sessionID,
			})
		c.Set("context", ctx)
		clientIP := c.ClientIP()
		clientUserAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		entryLog := dLogger.Logger.WithFields(logrus.Fields{
			"hostname":  dLogger.Lops.Hostname,
			"clientIP":  clientIP,
			"method":    c.Request.Method,
			"path":      path,
			"caller":    referer,
			"userAgent": clientUserAgent,
			"requestID": requestID,
			"sessionID": sessionID,
		})
		entryLog.Info("Request Initiated")
		c.Next()
		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		statusCode := c.Writer.Status()
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}

		entry := dLogger.Logger.WithFields(logrus.Fields{
			"hostname":   dLogger.Lops.Hostname,
			"statusCode": statusCode,
			"latency":    latency, // time to process
			"clientIP":   clientIP,
			"method":     c.Request.Method,
			"path":       path,
			"host":       c.Request.Host,
			"caller":     referer,
			"dataLength": dataLength,
			"userAgent":  clientUserAgent,
			"requestID":  requestID,
			"sessionID":  sessionID,
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			msg := "Request Completed"
			if statusCode > 499 {
				entry.Error(msg)
			} else if statusCode > 399 {
				entry.Warn(msg)
			} else {
				entry.Info(msg)
			}
		}
	}
}
