package sam

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSQLSplit(t *testing.T) {
	cases := []struct {
		data   string
		expect []string
	}{
		{
			data:   "a;",
			expect: []string{"a;"},
		},
		{
			data:   "a;\nb;",
			expect: []string{"a;", "b;"},
		},
		{
			data:   "a\nb;",
			expect: []string{"a\nb;"},
		},
		{
			data:   "a;  \nb;",
			expect: []string{"a;", "b;"},
		},
		{
			data:   "a;  \nb",
			expect: []string{"a;", "b"},
		},
		{
			data:   "a;  \nb\n\n\n\n",
			expect: []string{"a;", "b"},
		},
	}

	for x, c := range cases {
		t.Run(fmt.Sprintf("case#%d", x+1), func(t *testing.T) {
			ret := split([]byte(c.data))
			actual := make([]string, len(ret))
			for x, l := range ret {
				actual[x] = string(l)
			}

			if reflect.DeepEqual(actual, c.expect) {
				return
			}

			t.Errorf("expected %#v, got %#v", c.expect, actual)
		})
	}
}
