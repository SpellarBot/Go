// MySQL Client
// @Author: Golion
// @Date: 2017.3

package mysql

import (
	"fmt"
	"time"
	"sync/atomic"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLClient struct {
	ConnStr     string
	maxIdleConn int
	maxConn 	int
	Debug       bool
	AutoUseDB   string       // [坑！请务必注意！]使用TCPHA的话，需要每次请求前USE DATABASE一次，要不然UCDCProxy可能会串库
	Logger      func(string)

	conn        *sql.DB
	ticker      *time.Ticker
	qps         int32
	qpsCnt      int32
}

func (m *MySQLClient) Init() error {
	if m.Logger == nil {
		m.Logger = func(msg string) {
			fmt.Printf(msg)
		}
	}
	if m.conn == nil {
		var err error
		m.conn, err = sql.Open("mysql", m.ConnStr)
		if m.maxIdleConn == 0 {
			m.maxIdleConn = 5
		}
		m.conn.SetMaxIdleConns(m.maxIdleConn)
		if m.maxConn == 0 {
			m.maxConn = 200
		}
		m.conn.SetMaxOpenConns(m.maxConn)
		if err != nil {
			m.Logger(fmt.Sprintf("[MySQLClient][Init] Error! Open MySQL Connection Failed! error=[%v]", err.Error()))
			return fmt.Errorf("[MySQLClient][Init] error=[%v]", err.Error())
		}
	}
	if m.ticker == nil {
		go m.countQPS()
	}
	return nil
}

func (m *MySQLClient) Close() {
	if m.conn != nil {
		m.conn.Close()
		m.conn = nil
	}
}

func (m *MySQLClient) UseDB(DBName string) (bool, error) {
	if m.conn == nil {
		return false, fmt.Errorf("[MySQLClient][UseDB] Error! Connect to MySQL Failed! Have You Init()?")
	}
	_, err := m.conn.Exec("USE " + Q(DBName))
	if err != nil {
		return false, fmt.Errorf("[MySQLClient][UseDB] error=[%v]", err.Error())
	}
	return true, nil
}

func (m *MySQLClient) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if m.conn == nil {
		return nil, fmt.Errorf("[MySQLClient][Query] Error! Connect to MySQL Failed! Have You Init()?")
	}
	atomic.AddInt32(&m.qpsCnt, 1)
	if m.AutoUseDB != "" {
		if _, err := m.UseDB(m.AutoUseDB); err != nil {
			return nil, fmt.Errorf("[MySQLClient][Query] error=[%v]", err.Error())
		}
	}
	rows, err := m.conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("[MySQLClient][Query] error=[%v]", err.Error())
	}
	if err = rows.Err(); err != nil {
		rows.Close()
		return nil, fmt.Errorf("[MySQLClient][Query] error=[%v]", err.Error())
	}
	return rows, nil
}

func (m *MySQLClient) GetOne(query string, result ...interface{}) (bool, error) {
	rows, err := m.Query(query)
	if err != nil {
		return false, fmt.Errorf("[MySQLClient][GetOne] error=[%v]", err.Error())
	}
	defer rows.Close()
	if err = rows.Err(); err != nil {
		return false, fmt.Errorf("[MySQLClient][GetOne] error=[%v]", err.Error())
	}
	if rows.Next() {
		err := rows.Scan(result...)
		if err != nil {
			return false, fmt.Errorf("[MySQLClient][GetOne] error=[%v]", err.Error())
		}
		return true, nil
	}
	if err = rows.Err(); err != nil {
		return false, fmt.Errorf("[MySQLClient][GetOne] error=[%v]", err.Error())
	}
	return false, fmt.Errorf("[MySQLClient][GetOne] Empty!")
}

func (m *MySQLClient) ExecAndGetResult(query string, args ...interface{}) (*sql.Result, error) {
	if m.conn == nil {
		return nil, fmt.Errorf("[MySQLClient][ExecAndGetResult] Error! Connect to MySQL Failed! Have You Init()?")
	}
	atomic.AddInt32(&m.qpsCnt, 1)
	if m.AutoUseDB != "" {
		if _, err := m.UseDB(m.AutoUseDB); err != nil {
			return nil, fmt.Errorf("[MySQLClient][ExecAndGetResult] error=[%v]", err.Error())
		}
	}
	res, err := m.conn.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("[MySQLClient][ExecAndGetResult] error=[%v]", err.Error())
	}
	return &res, nil
}

func (m *MySQLClient) Exec(query string, args ...interface{}) (bool, error) {
	if _, err := m.ExecAndGetResult(query, args...); err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (m *MySQLClient) ExecAndGetLastInsertId(query string, args ...interface{}) (int64, error) {
	if res, err := m.ExecAndGetResult(query, args...); err != nil {
		return -1, err
	} else {
		lastId, err := (*res).LastInsertId()
		if err != nil {
			return -1, fmt.Errorf("[MySQLClient][ExecAndGetLastInsertId] error=[%v]", err.Error())
		}
		return lastId, nil
	}
}

func (m *MySQLClient) ExecAndGetRowsAffected(query string, args ...interface{}) (int64, error) {
	if res, err := m.ExecAndGetResult(query, args...); err != nil {
		return -1, err
	} else {
		rowsAffected, err := (*res).RowsAffected()
		if err != nil {
			return -1, fmt.Errorf("[MySQLClient][ExecAndGetRowsAffected] error=[%v]", err.Error())
		}
		return rowsAffected, nil
	}
}

// Equals mysql_escape_string()
func Q(s string) string {
	var j int = 0
	if len(s) == 0 {
		return ""
	}
	tempStr := s[:]
	desc := make([]byte, len(tempStr)*2)
	for i := 0; i < len(tempStr); i++ {
		flag := false
		var escape byte
		switch tempStr[i] {
		case '\r':
			flag = true
			escape = '\r'
			break
		case '\n':
			flag = true
			escape = '\n'
			break
		case '\\':
			flag = true
			escape = '\\'
			break
		case '\'':
			flag = true
			escape = '\''
			break
		case '"':
			flag = true
			escape = '"'
			break
		case '\032':
			flag = true
			escape = 'Z'
			break
		default:
		}
		if flag {
			desc[j] = '\\'
			desc[j+1] = escape
			j = j + 2
		} else {
			desc[j] = tempStr[i]
			j = j + 1
		}
	}
	return string(desc[0:j])
}

func (m *MySQLClient) QPS() int32 {
	return m.qps
}

func (m *MySQLClient) countQPS() {
	m.ticker = time.NewTicker(1 * time.Second)
	for _ = range m.ticker.C {
		m.qps = m.qpsCnt
		m.qpsCnt = 0
	}
}