package visitor

import (
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/model"
	"fmt"
	"log"
	"reflect"
)

type AVisitor struct {
	tbls     map[string]Table
	joinType int
}

func NewAVisitor() AVisitor {
	return AVisitor{
		tbls: make(map[string]Table),
	}
}

func (v *AVisitor) getTblNames() string {
	l := TableList{}
	for _, t := range v.tbls {
		l = l.Append(t)
	}
	return l.String()
}

func (v *AVisitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	log.Println("Enter: ", reflect.TypeOf(in))
	switch n := in.(type) {
	case *ast.SelectStmt:
		v.parseSelectStmt(n, v.tbls)
		return in,true
	default:

	}
	return in, false
}

func (v *AVisitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	//log.Println("Leave: ", reflect.TypeOf(in))
	return in, true
}

func (v *AVisitor) parseSelectStmt(n *ast.SelectStmt, m map[string]Table) {
	tr := n.From.TableRefs

	v.joinType = int(tr.Tp)

	//log.Println(reflect.TypeOf(tr.Left))
	//log.Println(reflect.TypeOf(tr.Right))
	v.extractTblName(tr.Left, m)
	v.extractTblName(tr.Right, m)
}

func (v *AVisitor) extractTblName(i ast.ResultSetNode, m map[string]Table) {
	//log.Println("ex: ", reflect.TypeOf(i))
	if i != nil {
		switch l := i.(type) {
		case *ast.TableSource:
			//l.AsName
			//log.Println(reflect.TypeOf(l.Source))
			//v.tbls = append(v.tbls, l.Source.(*ast.TableName).Name.L)
			switch s := l.Source.(type) {
			case *ast.SelectStmt:
				name := model.NewCIStr(fmt.Sprintf("%s", l.AsName.L))
				log.Println("asname: ",name.L)
				var tbl Table
				if _, ok := m[name.L]; !ok {
					log.Println("couldn't get tbl")
					tbl = Table{Name: name, Children: make(map[string]Table)}
					m[tbl.Name.L] = tbl
				}
				v.parseSelectStmt(s, tbl.Children)
			case *ast.TableName:
				name := model.NewCIStr(fmt.Sprintf("%s.%s", s.Schema.L, s.Name.L))
				if _, ok := m[name.L]; !ok {
					tbl := Table{Name: name, Children: make(map[string]Table)}
					m[tbl.Name.L] = tbl
				}
			}
		case *ast.Join:
			v.extractTblName(l.Left, m)
			v.extractTblName(l.Right, m)

		}
	}
}
