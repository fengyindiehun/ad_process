package main

import (
	"database/sql"
	"fmt"
	"github.com/go-ini/ini"
	_ "github.com/go-sql-driver/mysql"
)

type mysqlProxy struct {
	username, password, host, port, dbname, tablefeed, tablecust string
	db                                                           *sql.DB
	isConnected                                                  bool
}

func (proxy *mysqlProxy) doConnect() bool {
	if proxy.isConnected == true {
		return true
	}

	var err error
	proxy.db, err = sql.Open("mysql", proxy.username+
		":"+proxy.password+
		"@tcp("+proxy.host+":"+
		proxy.port+")/"+proxy.dbname)
	if err != nil {
		proxy.isConnected = false
		return false
	} else {
		proxy.isConnected = true
		return true
	}
}

func (proxy *mysqlProxy) updateFeedInfo(feedId, custId, text, creativeType, tag string) error {
	sql := "INSERT INTO " +
		proxy.tablefeed +
		"(feedId, custId, text, creativeType, tag) values(?,?,?,?,?)" +
		" ON DUPLICATE KEY UPDATE feedId=?, custId=?, text=?, creativeType=?, tag=?"
	_, err := proxy.db.Exec(sql, feedId, custId, text, creativeType, tag,
		feedId, custId, text, creativeType, tag)
	return err
}

func (proxy *mysqlProxy) updateCustInfo(custId, custType string) error {
	sql := "INSERT INTO " +
		proxy.tablecust +
		"(custId, custType) values(?,?)" +
		" ON DUPLICATE KEY UPDATE custId=?, custType=?"
	_, err := proxy.db.Exec(sql, custId, custType, custId, custType)
	return err
}

func newMysqlProxy(section *ini.Section) (*mysqlProxy, error) {
	username := section.Key(MYSQL_SECTION_USERNAME).String()
	if username == "" {
		return nil, fmt.Errorf("mysql username parse error")
	}

	password := section.Key(MYSQL_SECTION_PASSWORD).String()
	if password == "" {
		return nil, fmt.Errorf("mysql password parse error")
	}

	host := section.Key(MYSQL_SECTION_HOST).String()
	if host == "" {
		return nil, fmt.Errorf("mysql host parse error")
	}

	port := section.Key(MYSQL_SECTION_PORT).String()
	if port == "" {
		return nil, fmt.Errorf("mysql port parse error")
	}

	dbname := section.Key(MYSQL_SECTION_DBNAME).String()
	if dbname == "" {
		return nil, fmt.Errorf("mysql dbname parse error")
	}

	tablefeed := section.Key(MYSQL_SECTION_TABLE_FEED).String()
	if tablefeed == "" {
		return nil, fmt.Errorf("mysql tablefeed parse error")
	}

	tablecust := section.Key(MYSQL_SECTION_TABLE_CUST).String()
	if tablecust == "" {
		return nil, fmt.Errorf("mysql tablecust parse error")
	}

	db, err := sql.Open("mysql", username+":"+password+"@tcp("+host+":"+port+")/"+dbname)
	if err != nil {
		return nil, err
	}

	return &mysqlProxy{
		username:    username,
		password:    password,
		host:        host,
		port:        port,
		dbname:      dbname,
		tablefeed:   tablefeed,
		tablecust:   tablecust,
		db:          db,
		isConnected: true,
	}, nil

	return nil, nil
}
