package visitor

import (
	"github.com/pingcap/parser/ast"
	"log"
	"reflect"
)

type AVisitor struct {
	tbls     []string
	joinType int
	children []AVisitor
}

func (v *AVisitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	log.Println("Enter: ", reflect.TypeOf(in))
	switch n := in.(type) {
	case *ast.SelectStmt:
		v.parseSelectStmt(n)
		return in, false
	default:

	}
	return in, false
}

func (v *AVisitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	log.Println("Leave: ", reflect.TypeOf(in))
	return in, true
}

func (v *AVisitor) parseSelectStmt(n *ast.SelectStmt) {
	tr := n.From.TableRefs

	v.joinType = int(tr.Tp)

	log.Println(reflect.TypeOf(tr.Left))
	log.Println(reflect.TypeOf(tr.Right))
	v.extractTblName(tr.Left)
	v.extractTblName(tr.Right)
}

func (v *AVisitor) extractTblName(i ast.ResultSetNode) {
	log.Println("ex: ", reflect.TypeOf(i))
	if i != nil {
		switch l := i.(type) {
		case *ast.TableSource:
			//l.AsName
			v.tbls = append(v.tbls, l.Source.(*ast.TableName).Name.L)
		case *ast.Join:
			v.extractTblName(l.Left)
			v.extractTblName(l.Right)

		}
	}
}
