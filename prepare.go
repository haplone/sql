package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:4000)/test")
	if err != nil {
		panic(err)
		return
	}
	defer db.Close()

	stmt, _ := db.Prepare(`INSERT INTO user (name, age) VALUES (?, ?)`)
	defer stmt.Close()

	str :=`
''
"",
"",
`
	_, err = stmt.Exec(str, 23)

	if err != nil {
		fmt.Printf("insert data error: %v\n", err)
		return
	}

}
