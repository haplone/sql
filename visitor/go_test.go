package visitor

import (
	"testing"
	"github.com/pingcap/parser"
	_ "github.com/pingcap/tidb/types/parser_driver"

	"os"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
)

type TblCheck struct {
	sql    string
	expect AVisitor
}

func TestAVisitor_Enter(t *testing.T) {
	sqls := []struct {
		sql    string
		expect AVisitor
	}{
		TblCheck{
			"insert into target select sum(a) from source",
			AVisitor{tbls: []string{"source"}},
		},
		TblCheck{
			"insert into target select sum(s1.a) from source s1, source2 s2 where s1.id = s2.id",
			AVisitor{tbls: []string{"source", "source2"}},
		},
		TblCheck{
			"insert into target select sum(s1.a) from source s1 join  source2 s2 on s1.id = s2.id",
			AVisitor{tbls: []string{"source", "source2"}},
		},
		TblCheck{
			"insert into target select sum(s1.a) from source s1, source2 s2,source3 s3 " +
				"where s1.id = s2.id and s2.id = s3.id",
			AVisitor{tbls: []string{"source", "source2", "source3"}},
		},
	}

	for _, c := range sqls {
		p := parser.New()
		ast, err := p.ParseOneStmt(c.sql, "", "")
		Check(err)

		v := AVisitor{}
		ast.Accept(&v)

		t.Log(v.tbls)
		t.Log(ast)
		assert.Equal(t, c.expect.tbls, v.tbls)
	}

}

func Check(err error) {
	if err != nil {
		//log.Println(err)
		panic(err)
	}
}

func getSql() string {
	f, err := os.Open("/data/example.sql")
	Check(err)
	c, err := ioutil.ReadAll(f)
	Check(err)
	return string(c)
}
