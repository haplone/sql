package visitor

import (
	"strings"
	"github.com/pingcap/parser/model"
	"sort"
	"fmt"
	"log"
)

type Table struct {
	Name     model.CIStr
	Children map[string]Table
}

func (t *Table) String() string {
	if len(t.Children) > 0 {
		var tl TableList
		for _, tbl := range t.Children {
			tl = tl.Append(tbl)
		}
		return fmt.Sprintf("%s[%s]", t.Name.L, tl.String())
	}
	return t.Name.L
}

type TableList []*Table

func (t TableList) Len() int {
	return len(t)
}

func (t TableList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t TableList) Less(i, j int) bool {
	return strings.Compare(t[i].Name.L, t[j].Name.L) < 0
}

func (t TableList) String() string {
	var sl []string

	sort.Sort(t)
	for _, i := range t {
		sl = append(sl, i.String())
	}
	return strings.Join(sl, ",")
}

func (t TableList) Append(tbl Table) TableList {
	t = append(t, &tbl)
	return t
}

type SelectStmt struct {
	Nodes []SelectNode
}

func (s *SelectStmt) Append(node SelectNode) {
	//log.Println(reflect.TypeOf(node) , node.String())
	s.Nodes = append(s.Nodes, node)
}

func (s *SelectStmt) String() string {
	if s == nil || s.Nodes == nil {
		log.Println("has no stmt")
		return ""
	}
	skip := make(map[int]bool)
	var r []string
	for idx, i := range s.Nodes {
		if _, ok := skip[idx]; !ok {
			switch n := i.(type) {
			case *SJoin:
				r = append(r, n.Type+" join "+s.Nodes[idx+1].String()+" on "+n.Text())
				skip[idx+1] = true
			default:
				r = append(r, i.String())
			}
		}
	}
	return "select from " + strings.Join(r, " ")
}

type SelectNode interface {
	String() string
}
type STable struct {
	//SelectNode
	Name   model.CIStr
	Alia   model.CIStr
	Schema model.CIStr
}

func (st *STable) String() string {
	//log.Println(fmt.Sprintf("%s.%s", st.Schema.L, st.Name.L))
	return fmt.Sprintf("%s.%s %s", st.Schema.L, st.Name.L, st.Alia.L)
}

type STmpTable struct {
	Alia model.CIStr
	//SelectNode
	Stmt *SelectStmt
}

func (stt *STmpTable) String() string {
	//log.Println(stt.Alia.L)
	return fmt.Sprintf("(%s) %s ", stt.Stmt.String(), stt.Alia.L)
	//return stt.Alia.L
}

type SJoin struct {
	Type  string
	Nodes []ColumnExpr
}

func (j *SJoin) Text() string {
	var r [] string

	for _, n := range j.Nodes {
		r = append(r, n.String())
	}
	return strings.Join(r, " ")
}

func (j *SJoin) String() string {
	var r [] string
	r = append(r, j.Type)

	for _, n := range j.Nodes {
		r = append(r, n.String())
	}
	return "====" + strings.Join(r, " ") + "===="
}

type ColumnExpr struct {
	L  ColumnNode
	Op string
	R  ColumnNode
}

func (ce *ColumnExpr) String() string {
	return fmt.Sprintf("%s %s %s", ce.L.String(), ce.Op, ce.R.String())
}

type ColumnNode struct {
	C      string
	Schema model.CIStr
	Table  model.CIStr
	Name   model.CIStr
}

func (c *ColumnNode) String() string {
	return fmt.Sprintf("%s.%s.%s", c.Schema.L, c.Table.L, c.Name.L)
}
