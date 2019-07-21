package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ltltlt/ip-recorder/server/utils"
)

var (
	mysqlUser string
	mysqlPass string
	mysqlHost = "mysql_host"
	dbName    = "statistics"
	table     = "ip"
	/*
		+-------+------------------------------------------------------+
		| Table | Create Table                                         |
		+-------+------------------------------------------------------+
		| ip    | CREATE TABLE `ip` (                                  |
		|       |   `ip` varchar(16) NOT NULL,                         |
		|       |   `time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP |
		|       | ) ENGINE=InnoDB DEFAULT CHARSET=latin1               |
		+-------+------------------------------------------------------+

	*/
	createTableStatement = "create table if not exists ip (ip varchar(130) not null, time datetime not null default current_timestamp)"
	insertStatement      = fmt.Sprintf("insert into %s(ip, time) values(?, ?)", table)

	db *sql.DB
)

func init() {
	var ok1, ok2 bool
	mysqlUser, ok1 = os.LookupEnv("MYSQL_USER")
	mysqlPass, ok2 = os.LookupEnv("MYSQL_PASS")

	if !ok1 && !ok2 {
		log.Printf("you may forget to setup MYSQL_USER or MYSQL_PASS")
	}
}

// late init, when two container starts, go can start just after mysql starts
// in this case, it will be a problem when go tries to execute create table and mysql still
// not fully initialized
var once sync.Once

func beforeWrite() {
	once.Do(func() {
		var err error
		db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/", mysqlUser, mysqlPass, mysqlHost))
		if err != nil {
			log.Panicf("fail initialize mysql: %v", err)
		}
		_, err = db.Exec("create database if not exists " + dbName)
		if err != nil {
			log.Panicf("fail to initialize mysql by create database: %v", err)
		}
		_, err = db.Exec("use " + dbName)
		if err != nil {
			log.Panicf("fail to initialize mysql by switch db: %v", err)
		}
		_, err = db.Exec(createTableStatement)
		if err != nil {
			log.Panicf("fail to initialize mysql by create table ip: %v", err)
		}
	})
}

func AccessLog(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := utils.GetIP(r)
		ctx, _ := context.WithTimeout(context.Background(), time.Second*3)

		beforeWrite()

		result, err := db.ExecContext(ctx, insertStatement, ip, time.Now())
		if err != nil {
			log.Printf("fail to add access log: %v", err)
		} else {
			id, err := result.LastInsertId()
			if err != nil {
				log.Printf("fail to get last insert id: %v", err)
			}
			rowAffected, err := result.RowsAffected()
			if err != nil {
				log.Printf("fail to get rowAffected: %v", err)
			}
			log.Printf("add access log: %v(id), %v(affected row num)", id, rowAffected)
		}

		handler(w, r)
	}
}
