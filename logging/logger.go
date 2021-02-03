package logging

import (
	"bytes"
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
	file = strings.Trim(file, " ")
	funcname = strings.Trim(funcname, " ")
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
	if _, ok := fields["caller"]; !ok {
		fields["caller"] = fmt.Sprintf("%s:%d\n%s", file, line, funcname)
	}
	fields["appID"] = dLogger.Lops.APPID
	fields["requestID"] = corRelationID["requestID"]
	fields["sessionID"] = corRelationID["sessionID"]
	entry := dLogger.Logger.WithFields(fields)
	entry.Log(cb, MessageKey)
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
		fields := logrus.Fields{
			"referer":   c.Request.Referer(),
			"clientIP":  c.ClientIP(),
			"host":      c.Request.Host,
			"hostname":  dLogger.Lops.Hostname,
			"method":    c.Request.Method,
			"path":      c.Request.URL.Path,
			"query":     c.Request.URL.RawQuery,
			"requestID": requestID,
			"sessionID": sessionID,
			"userAgent": c.Request.UserAgent(),
		}

		var level = new(logrus.Level)
		*level = logrus.InfoLevel

		//deferred request log
		defer dLogger.WriteLogs(ctx, fields, *level, "access-log")
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
			fields["requestBody"] = string(bodyBytes)
		}
		// Restore the io.ReadCloser to its original state
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		fields["statusCode"] = 0
		c.Next()
		stop := time.Since(start)
		fields["latency"] = int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		code := c.Writer.Status()

		fields["statusCode"] = code
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}
		fields["dataLength"] = dataLength

		if len(c.Errors) > 0 {
			fields["error"] = c.Errors.ByType(gin.ErrorTypePrivate).String()
			*level = logrus.ErrorLevel
		} else if code > 499 {
			*level = logrus.ErrorLevel
		} else if code > 399 {
			*level = logrus.WarnLevel
		}
	}
}
