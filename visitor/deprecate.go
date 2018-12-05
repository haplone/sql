package visitor

import (
	"log"
	"reflect"
	"strings"
	"fmt"
	"github.com/pingcap/parser/ast"
)

type FVisitor struct {
}

func (v *FVisitor) Enter(in ast.Node) (node ast.Node, skipChildren bool) {
	log.Printf("Enter: %s \r\n", reflect.TypeOf(in))
	return in, false
}

func (v *FVisitor) Leave(in ast.Node) (node ast.Node, ok bool) {
	log.Printf("Leave: %s\r\n", reflect.TypeOf(in))
	return in, true
}

type TblNameVisitor struct {
}

func (v *TblNameVisitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	log.Println("tp: %s", reflect.TypeOf(in))
	switch a := in.(type) {
	case *ast.TableName:
		log.Printf("tblName: %s.%s", a.Schema.L, a.Name.L)
	case *ast.ColumnName:
		log.Printf("colName: %s.%s %s", a.Schema.L, a.Table.L, a.Name.L)
	case *ast.ColumnNameExpr:
		log.Printf("colNameExpr: %v", a)
		String(a)
	case *ast.FieldList:
		String(a)
	default:
	}

	return in, false
}

func (v *TblNameVisitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

func String(in ast.Node) string {
	switch a := in.(type) {
	case *ast.ColumnNameExpr:
		r := a.Refer
		s := strings.Builder{}
		s.WriteString(a.Name.Name.L)
		if r != nil && r.Referenced {
			s.WriteString(fmt.Sprintf("%s - %s", r.TableName.Name.L, reflect.TypeOf(r.Expr)))
		}

		log.Printf("colNameExpStr: %s", s.String())
		return s.String()
	case *ast.FieldList:
		for _, f := range a.Fields {
			if f.WildCard != nil {
				log.Printf("fields: %s ,%s.%s  %s", f.AsName.L, f.WildCard.Schema.L, f.WildCard.Table.L, f.WildCard.Text())
			}
			if f.Expr != nil {
				switch e := f.Expr.(type) {
				case *ast.AggregateFuncExpr:
					log.Printf("AggregateFuncExpr: %v, %s %v", e.Distinct, e.F, e.Type.Tp)
					for _, a := range e.Args {
						log.Printf("tp: %s, ",reflect.TypeOf(a))
						switch n := a.(type){
						case *ast.ColumnNameExpr:
							cn := n.Name
							log.Printf("afe source: %s.%s.%s",cn.Schema.L,cn.Table.L,cn.Name.L)
						}
						//String(a)
					}
				}
				log.Println("========", reflect.TypeOf(f.Expr))
			}
		}
	}
	return ""
}

