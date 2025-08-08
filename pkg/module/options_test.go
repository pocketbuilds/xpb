package module

import (
	"reflect"
	"testing"
)

func TestReCliArg(t *testing.T) {
	type test struct {
		name   string
		arg    string
		expect []string
	}

	scenarios := []test{
		{
			name: "full arg",
			arg:  "github.com/fake/plugin@v0.0.1=my/replacement",
			expect: []string{
				"github.com/fake/plugin",
				"v0.0.1",
				"my/replacement",
			},
		},
		{
			name: "only module",
			arg:  "github.com/fake/plugin",
			expect: []string{
				"github.com/fake/plugin",
				"",
				"",
			},
		},
		{
			name: "only module and version",
			arg:  "github.com/fake/plugin@v0.0.1",
			expect: []string{
				"github.com/fake/plugin",
				"v0.0.1",
				"",
			},
		},
		{
			name: "only module and replacement",
			arg:  "github.com/fake/plugin=my/replacement",
			expect: []string{
				"github.com/fake/plugin",
				"",
				"my/replacement",
			},
		},
	}

	for _, s := range scenarios {
		result := reCliArg.FindStringSubmatch(s.arg)
		result = result[1:]
		if !reflect.DeepEqual(result, s.expect) {
			t.Errorf("%s failed: got: %#v expected %#v\n", s.name, result, s.expect)
		}
	}
}
