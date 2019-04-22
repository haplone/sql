package main

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// -uroot -h127.0.0.1 -P4000 -up10080
func main() {
	dbUser := flag.String("u", "root2", "user for database")
	dbPwd := flag.String("p", "", "password for database")
	dbHost := flag.String("h", "127.0.0.1", "host for database")
	dbPort := flag.String("P", "4000", "port for database")
	diskUsagePort := flag.String("up", "10080", "port for tidb disk usage")
	batchSize := flag.Int64("bs", 1000, "per insert sql contains how much tbl")
	flag.Parse()

	cfg := fmt.Sprintf("%s:%s@tcp(%s:%s)/test", *dbUser, *dbPwd, *dbHost, *dbPort)
	log.Printf("sql client url: %s", cfg)
	t := NewTask(cfg, *diskUsagePort, *dbHost, *batchSize)
	defer t.DBClient.Close()
	t.ParseDbs()
	t.ParseTbls()
	t.FillTblUsage()
	t.Tbl2Sql()
	//for _, db := range t.Dbs {
	//	log.Println("===============")
	//	log.Printf("db  : %s \n", db.Name)
	//	for _, tbl := range db.Tbls {
	//		log.Printf("tbl: %s\n", tbl.Name)
	//	}
	//}
}

func (t *DiskUsageTask) ParseDbs() {
	rds, err := t.DBClient.Query("show databases")
	check2(err)
	defer rds.Close()

	for rds.Next() {
		var n string
		rds.Scan(&n)
		log.Printf("--- db name: %s", n)
		db := Db{Name: n}
		t.Dbs = append(t.Dbs, &db)
		log.Printf("--- db name: %s", db.Name)
	}
}

func (t *DiskUsageTask) ParseTbls() {
	for _, db := range t.Dbs {
		sql := fmt.Sprintf("use %s", db.Name)
		log.Printf("--- sql: %s", sql)
		t.DBClient.Exec(sql)
		rds, err := t.DBClient.Query("show tables")
		check2(err)
		for rds.Next() {
			var n string
			rds.Scan(&n)
			log.Printf("--- tbl name: %s", n)
			db.Tbls = append(db.Tbls, &Tbl{Name: n, DbName: db.Name})
		}
	}
}

func (t *DiskUsageTask) FillTblUsage() {
	for _, db := range t.Dbs {
		for _, tbl := range db.Tbls {
			url := fmt.Sprintf("http://%s:%s/tables/%s/%s/disk-usage", t.DbHost, t.DiskUsagePort, db.Name, tbl.Name)

			resp, err := http.Get(url)
			check2(err)
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			check2(err)

			bytesBuffer := bytes.NewBuffer(body)
			binary.Read(bytesBuffer, binary.BigEndian, &tbl.Usage)

			log.Printf("%s.%s : %d", tbl.DbName, tbl.Name, tbl.Usage)
		}
	}
}

func (t *DiskUsageTask) Tbl2Sql() {
	var count int64
	logDay := time.Now().Format("2006-01-02")
	sqlPrefix := "insert into tidb_monitor.table_disk_usage (db,tbl,disk_usage,log_day) values "
	var sqls []string

	for _, db := range t.Dbs {
		for _, tbl := range db.Tbls {
			//t.TblChan <- tbl.clone()
			count += 1

			sqls = append(sqls, fmt.Sprintf("('%s','%s',%d,'%s')", tbl.DbName, tbl.Name, tbl.Usage, logDay))

			if count%t.BatchSize == 0 {
				//for ttt := range t.TblChan {
				//}
				sql := sqlPrefix + strings.Join(sqls, ",")
				log.Println(sql)
				_, err := t.DBClient.Exec(sql)
				check2(err)
				sqls = make([]string, 0)

			}
		}
	}
	if len(sqls) > 0 {
		sql := sqlPrefix + strings.Join(sqls, ",")
		log.Println(sql)
		_, err := t.DBClient.Exec(sql)
		check2(err)
	}
	log.Printf("-- we got %d tbls", count)
}

func NewTask(cfg string, diskUsagePort string, dbHost string, batchSize int64) *DiskUsageTask {
	client := GetMysql(cfg)
	return &DiskUsageTask{
		DbSource:      cfg,
		DBClient:      client,
		TblChan:       make(chan Tbl, 10000),
		DiskUsagePort: diskUsagePort,
		DbHost:        dbHost,
		BatchSize:     batchSize,
		//Dbs:      make([]*Db, 1),
	}
}

type DiskUsageTask struct {
	DbSource string
	DBClient *sql.DB
	Dbs      []*Db
	TblChan  chan Tbl

	DiskUsagePort string
	DbHost        string
	BatchSize     int64
}

type Db struct {
	Name string
	Tbls []*Tbl
}

type Tbl struct {
	Name   string
	DbName string
	Usage  uint32
}

func (tbl *Tbl) clone() Tbl {
	return Tbl{
		Name:   tbl.Name,
		DbName: tbl.DbName,
		Usage:  tbl.Usage,
	}
}

func GetMysql(cfg string) *sql.DB {
	db, err := sql.Open("mysql", cfg)
	if err != nil {
		log.Fatal("can not connect tidb ")
	}

	return db
}

func check2(err error) {
	if err != nil {
		log.Println(err)
	}
}
