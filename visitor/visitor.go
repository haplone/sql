package visitor

import "github.com/pingcap/parser/ast"

type AVisitor struct {
}

func (v *AVisitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	return in, false
}

func (v *AVisitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}
