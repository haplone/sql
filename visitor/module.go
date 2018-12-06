package visitor

import (
	"strings"
	"github.com/pingcap/parser/model"
	"sort"
	"fmt"
)

type Table struct {
	Name model.CIStr
	Children map[string]Table
}

func (t *Table) String() string {
	if len(t.Children) >0 {
		var tl TableList
		for _,tbl := range t.Children{
			tl = tl.Append(tbl)
		}
		return fmt.Sprintf("%s[%s]",t.Name.L,tl.String())
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

func (t TableList) Append(tbl Table) TableList{
	t = append(t, &tbl)
	return t
}