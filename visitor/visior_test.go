package visitor

import (
	"testing"
	"github.com/pingcap/parser"
	_ "github.com/pingcap/tidb/types/parser_driver"

	"os"
	"io/ioutil"
	//"github.com/stretchr/testify/assert"
	"log"
)

// toDo analysis union | union all | exist (select from) | not in
//func TestAVisitor_Enter_simple(t *testing.T) {
//	sqls := []struct {
//		sql    string
//		expect string
//	}{
//		{
//			"insert into target select sum(a) from s.source",
//			"s.source",
//		},
//		{
//			"insert into target select sum(s1.a) from s.source s1, s.source2 s2 where s1.id = s2.id",
//			"s.source,s.source2",
//		},
//		{
//			"insert into target select sum(s1.a) from s.source s1 join  s.source2 s2 on s1.id = s2.id",
//			"s.source,s.source2",
//		},
//		{
//			"insert into target select sum(s1.a) from s.source s1, s.source2 s2,s.source3 s3 " +
//				"where s1.id = s2.id and s2.id = s3.id",
//			"s.source,s.source2,s.source3",
//		},
//		{
//			`insert into target select sum(s1.a) from s.source s1
//					join  s.source2 s2 on s1.id =s2.id
//					join s.source3 s3 on s2.id = s3.id`,
//			"s.source,s.source2,s.source3",
//		},
//	}
//
//	for _, c := range sqls {
//		p := parser.New()
//		ast, err := p.ParseOneStmt(c.sql, "", "")
//		Check(err)
//
//		v := NewAVisitor()
//		ast.Accept(&v)
//
//		//t.Log(v.tbls)
//		t.Log(ast)
//		assert.Equal(t, c.expect, v.getTblNames())
//	}
//
//}

func TestAVisitor_Enter_multi(t *testing.T) {
	sqls := []struct {
		sql    string
		expect string
	}{
		{`
		INSERT INTO target
		SELECT t1.tc1, t1.tc2, t1.tc3, t1.tc4, t1.tc5, t1.ttime
		FROM (SELECT t1.id AS tc1, t1.tc2, t1.tc3, t1.tc4, DATE_FORMAT(t1.tc5,'%Y-%m-%d') AS tc5,
					DATE_FORMAT(t1.ttime,'%Y-%m-%d') AS ttime, CONCAT('exm',t4.id) AS c1_old
				FROM ss.source1 t1
				LEFT JOIN ss.source2 t4 ON t1.tc4 = t4.id) t1`,
			"t1[ss.source1,ss.source2]",
		},
		{`
		INSERT INTO target
		SELECT t1.tc1, t1.tc2, t1.tc3, t1.tc4, t1.tc5, t1.ttime
		FROM (SELECT t1.id AS tc1, t1.tc2, t1.tc3, t1.tc4, DATE_FORMAT(t1.tc5,'%Y-%m-%d') AS tc5,
					DATE_FORMAT(t1.ttime,'%Y-%m-%d') AS ttime, CONCAT('exm',t4.id) AS c1_old
				FROM ss.source1 t1
				LEFT JOIN ss.source2 t4 ON t1.tc4 = t4.id) t1
				LEFT JOIN (SELECT DISTINCT c1 FROM ss.source3) t2 ON t2.c1 = t1.c1_new
				LEFT JOIN (SELECT DISTINCT tc1, c2 FROM ss.source4 WHERE c3 = 'haha'
							AND c2 IS NOT NULL AND c2 <> "null"
				) t3 ON t3.tc1 = t1.c1
		WHERE t1.ttime >= "2018-09-30"`,
			"t1[ss.source1,ss.source2],t2[ss.source3],t3[ss.source4]",
		},
	}

	r := make([]string,len(sqls))
	for idx, c := range sqls {
		p := parser.New()
		ast, err := p.ParseOneStmt(c.sql, "", "")
		Check(err)

		v := NewAVisitor()
		ast.Accept(&v)

		//t.Log(v.getTblNames())
		//t.Log(ast)

		r[idx] = v.getTblNames()
		//assert.Equal(t, c.expect, v.getTblNames())
	}

	for idx,c:= range sqls{
		log.Println("=========")
		log.Println(c.sql)
		log.Println(r[idx])
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
