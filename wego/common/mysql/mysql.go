package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlClient struct {
	Host        string
	Port        string
	Dbname      string
	Username    string
	Password    string
	MaxOpenConn int // 最大连接数
	MaxIdleConn int // 最大空闲数
	Qps         int32
	qpscnt      int32
	Db          *sql.DB
	Logger      func(interface{})
}

func (m *MysqlClient) connect() error {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", m.Username, m.Password, m.Host, m.Port, m.Dbname)
	db, err := sql.Open("mysql", connStr)
	m.Logger("[Mysql][Init]:" + connStr)
	if err == nil {
		m.Db = db
		if m.MaxIdleConn >= 0 && m.MaxOpenConn >= 0 {
			if m.MaxIdleConn >= m.MaxOpenConn {
				m.Db.SetMaxOpenConns(m.MaxOpenConn)
				m.Db.SetMaxIdleConns(m.MaxIdleConn)
			} else {
				m.Db.SetMaxOpenConns(m.MaxOpenConn)
			}
		} else {
			m.Db.SetMaxOpenConns(200)
		}
		_, err = m.Db.Exec(fmt.Sprintf("USE %s", m.Dbname))
		if err != nil {
			m.Db = nil
			m.Logger(fmt.Sprintf("[Mysql][InitError]: %s", err.Error()))
		} else {
			return nil
		}
	}
	m.Logger(fmt.Sprintf("[Mysql][InitError]: %s", err.Error()))
	return err
}

func (m *MysqlClient) Init() error {
	if m.Logger == nil {
		m.Logger = func(s interface{}) {
			fmt.Println(s)
		}
	}
	err := m.connect()
	go m.countQps()
	go m.keepAlive()
	return err
}

func (m *MysqlClient) QueryGetOne(query string, args ...interface{}) (result *sql.Row, err error) {
	if m.Db != nil {
		atomic.AddInt32(&(m.qpscnt), 1)
		result = m.Db.QueryRow(query, args...)
		return result, nil
	}
	m.Logger("[Mysql][QueryError]:Db is nil")
	return nil, errors.New("Db Is nil")
}

func (m *MysqlClient) QueryGetAll(query string, args ...interface{}) (result *sql.Rows, err error) {
	if m.Db != nil {
		atomic.AddInt32(&(m.qpscnt), 1)
		result, err = m.Db.Query(query, args...)
		if err == nil {
			if err = result.Err(); err == nil {
				return result, nil
			} else {
				m.Logger(fmt.Sprintf("[Mysql][QueryError]: %s", err.Error()))
			}
		}
		m.Logger(fmt.Sprintf("[Mysql][QueryError]: %s", err.Error()))
		return nil, err
	}
	m.Logger("[Mysql][QueryError]:Db is nil")
	return nil, errors.New("Db Is nil")
}

func (m *MysqlClient) ExecGetEffect(exec string, args ...interface{}) (result sql.Result, err error) {
	if m.Db != nil {
		atomic.AddInt32(&(m.qpscnt), 1)
		result, err = m.Db.Exec(exec, args...)
		if err != nil {
			m.Logger(fmt.Sprintf("[Mysql][ExecError]: %s", err.Error()))
		}
		return result, err
	} else {
		m.Logger("[Mysql][ExecError]:Db is nil")
		return result, errors.New("Db Is nil")
	}
}

func (m *MysqlClient) Close() {
	if m.Db != nil {
		m.Db.Close()
		m.Db = nil
	}
}

func (m *MysqlClient) keepAlive() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		<-ticker.C
		err := m.Db.Ping()
		if err != nil {
			m.connect()
			m.Logger(fmt.Sprintf("[Mysql][Died]:%s", err.Error()))
		}
	}
}

func (m *MysqlClient) countQps() {
	ticker := time.NewTicker(1 * time.Second)
	for _ = range ticker.C {
		m.Qps = m.qpscnt
		m.qpscnt = 0
	}
}
