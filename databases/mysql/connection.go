package mysql

import (
	"github.com/ziutek/mymysql/autorc"
	mymysql "github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/thrsafe" // Thread safe engine
	"github.com/scalia/mysynql/log"
)

type Connection struct {
	Conn *autorc.Conn
	DbName string
}

var conns map[string]Connection

func init() {
	conns = make(map[string] Connection)
}

func NewConnection(host, user, password, database string) Connection {
	dsn := user+":"+password+"@"+host+"/"+database
	conn, ok := conns[dsn]
	if ok {
		log.Debug("Reusing connection "+dsn)
		return conn
	}

	log.Debug("Creating new connection to "+dsn)

	db := autorc.New("tcp", "", host, user, password, database)
	db.Register("SET NAMES utf8")

	newconn := Connection{db, database}
	conns[dsn] = newconn

	return newconn
}

func (conn *Connection) Prepare(sql string) (*autorc.Stmt, error) {
	return conn.Conn.Prepare(sql)
}

func (conn *Connection) Query(sql string, args ...interface{}) ([]mymysql.Row, mymysql.Result, error) {
	return conn.Conn.Query(sql, args...)
}
