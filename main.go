package main

import (
	"github.com/haplone/sql/visitor"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	_ "github.com/pingcap/tidb/types/parser_driver"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	var sqlContent = "select cn1 from t1;select * from t2"
	sqlContent = getSql()
	//sqlContent = "insert into tt values('1')"
	sqlContent = "insert into db2.tt select a from db1.b"
	sqlContent = "insert into db2.tt select sum(a) aa from db1.b"

	p := parser.New()

	asts, _, err := p.Parse(sqlContent, "", "")

	check(err)
	var count int32 = 0
	for _, a := range asts {
		log.Println(a)

		switch tp := a.(type) {
		case *ast.InsertStmt:
			if count == 0 {

				log.Println(tp)
				f := visitor.NewAVisitor()
				tp.Accept(&f)
				log.Printf("======== target table: %s", tp.Table.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName).Name.L)
			}
			count += 1
		}
	}

}

func getSql() string {
	f, err := os.Open("/data/example.sql")
	check(err)
	c, err := ioutil.ReadAll(f)
	check(err)
	return string(c)
}

func check(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
