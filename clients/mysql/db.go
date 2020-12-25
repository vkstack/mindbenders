package mysql

import (
	"context"
	"database/sql"
	"time"

	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"
	"gitlab.com/dotpe/mindbenders/interfaces"
)

const (
	maxOpenConnection    = 10
	maxConnectionTimeout = time.Hour * 1
	maxIdleConnection    = 5

	//ErrDBConn ...
	ErrDBConn = "db_err_conn"
	//ErrDBRead ...
	ErrDBRead = "db_err_read"
	//ErrDBInsert ...
	ErrDBInsert = "db_err_insert"
	//ErrDBUpdate ...
	ErrDBUpdate = "db_err_update"
	//ErrDBWrite ...
	ErrDBWrite = "db_err_write"
	//ErrDBSQL ...
	ErrDBSQL = "db_err_write"
	//ErrDBTxn ...
	ErrDBTxn = "db_err_txn"
)

//Option ...
type Option struct {
	Host,
	Port,
	User,
	Pass,
	DB string

	MaxOpenConnection,
	MaxIdleConnection int
	MaxConnectionTimeOut time.Duration
}

func (ops *Option) setdefaults() {
	if ops.MaxConnectionTimeOut == 0 {
		ops.MaxConnectionTimeOut = maxConnectionTimeout
	}
	if ops.MaxOpenConnection == 0 {
		ops.MaxOpenConnection = maxOpenConnection
	}
	if ops.MaxIdleConnection == 0 {
		ops.MaxIdleConnection = maxIdleConnection
	}
}

//Init ...
func Init(option Option) (*sql.DB, error) {
	connectionstr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", option.User, option.Pass, option.Host, option.Port, option.DB)
	if option.User == "" || option.Pass == "" || option.Host == "" || option.Port == "" || option.DB == "" {
		return nil, fmt.Errorf("invalid connection detail: %s", connectionstr)
	}
	option.setdefaults()
	Client, err := sql.Open("mysql", connectionstr)
	Client.SetMaxOpenConns(option.MaxOpenConnection)
	Client.SetMaxIdleConns(option.MaxIdleConnection)
	Client.SetConnMaxLifetime(option.MaxConnectionTimeOut)
	if err != nil {
		return nil, err
	}
	fmt.Printf("startup:MySQL Connection established with: %s\n", connectionstr)
	return Client, nil
}

//WriteDBError ...
func WriteDBError(ctx context.Context, logger interfaces.ILogger, err error, db, etype string) {
	pc, file, line, _ := runtime.Caller(1)
	funcname := runtime.FuncForPC(pc).Name()
	logger.WriteLogs(ctx, logrus.Fields{
		"caller": fmt.Sprintf("%s:%d\n%s", file, line, funcname),
		"error":  error.Error(err),
	}, logrus.ErrorLevel, "MySQLError", etype, db)
}
