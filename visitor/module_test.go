package visitor

import (
	"testing"
	"github.com/pingcap/parser/model"
	"github.com/stretchr/testify/assert"
	"sort"
)

func TestTable_String(t *testing.T) {
	cl := []struct {
		in string
		out string
	}{
		{"",""},
		{"SDF","sdf"},
		{"你好","你好"},
	}
	for _,i := range cl{
		name := model.NewCIStr(i.in)
		tbl := Table{Name: name,Children: make(map[string]Table)}
		assert.Equal(t,i.out,tbl.String())
	}
}

func TestTableList(t *testing.T) {
	cl := []struct{
		in []string
		out string
	}{
		{[]string{"c","a","b"},"a,b,c"},
		{[]string{"a","a1","aa","ab"},"a,a1,aa,ab"},
	}

	for _,i := range cl{
		l := TableList{}
		for _,n := range i.in{
			l = l.Append(Table{Name: model.NewCIStr(n)})
		}
		sort.Sort(l)
		assert.Equal(t,i.out,l.String())
	}
}
