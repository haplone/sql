package visitor

import (
	"github.com/pingcap/parser/ast"
	"log"
	"reflect"
)

type AVisitor struct {
	Select *SelectStmt
}

func NewAVisitor() AVisitor {
	return AVisitor{
		Select: &SelectStmt{},
	}
}

func (v *AVisitor) getTblNames() string {
	//log.Println("getTblName ",v.Select.String())
	return v.Select.String()
}

func (v *AVisitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	log.Println("Enter: ", reflect.TypeOf(in))
	switch n := in.(type) {
	case *ast.SelectStmt:
		v.parseSelectStmt(n.From.TableRefs, v.Select)
		return in, true
	default:

	}
	return in, false
}

func (v *AVisitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	//log.Println("Leave: ", reflect.TypeOf(in))
	return in, true
}

func (v *AVisitor) parseSelectStmt(n *ast.Join, stmt *SelectStmt) {
	//log.Println(reflect.TypeOf(tr.Left))
	//log.Println(reflect.TypeOf(tr.Right))
	v.extractTblName(n.Left, stmt)

	if n.On != nil {
		log.Println(reflect.TypeOf(n.On.Expr))
		switch oe := n.On.Expr.(type) {
		case *ast.BinaryOperationExpr:
			log.Println(reflect.TypeOf(oe.L))
			lc := parseExprName(oe.L)
			log.Println(reflect.TypeOf(oe.Op))
			log.Println(reflect.TypeOf(oe.R))
			rc := parseExprName(oe.R)

			join := &SJoin{}
			switch n.Tp {
			case ast.CrossJoin:
				join.Type = "cross"
			case ast.LeftJoin:
				join.Type = "left"
			case ast.RightJoin:
				join.Type = "right"
			}
			join.Nodes = append(join.Nodes, ColumnExpr{
				L:  lc,
				Op: oe.Op.String(),
				R:  rc,
			})
			stmt.Append(join)
		}
	}
	v.extractTblName(n.Right, stmt)
}

func (v *AVisitor) extractTblName(i ast.ResultSetNode, stmt *SelectStmt) {
	//log.Println("ex: ", reflect.TypeOf(i))
	if i != nil {
		switch l := i.(type) {
		case *ast.TableSource:
			//l.AsName
			//log.Println(reflect.TypeOf(l.Source))
			//v.tbls = append(v.tbls, l.Source.(*ast.TableName).Name.L)
			switch s := l.Source.(type) {
			case *ast.SelectStmt:
				tmpTbl := &STmpTable{
					Alia: l.AsName,
					Stmt: &SelectStmt{},
				}
				stmt.Append(tmpTbl)
				v.parseSelectStmt(s.From.TableRefs, tmpTbl.Stmt)
			case *ast.TableName:
				st := &STable{
					Schema: s.Schema,
					Name:   s.Name,
					Alia:   l.AsName,
				}
				stmt.Append(st)
			}
		case *ast.Join:
			v.parseSelectStmt(l, stmt)
		}
	}
}

func (v *AVisitor) parseExpr(e ast.BinaryOperationExpr) {

}

func parseExprName(c ast.ExprNode) ColumnNode {
	switch e := c.(type) {
	case *ast.ColumnNameExpr:
		c := ColumnNode{
			Schema: e.Name.Schema,
			Table:  e.Name.Table,
			Name:   e.Name.Name,
		}
		return c
	}
	return ColumnNode{}
}
