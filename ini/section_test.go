package ini

import (
	"sort"
	"strings"
	"testing"
)

func TestSectionString(t *testing.T) {
	sec := newSection("sec1")
	sec.setProp("key1", "val1")
	sec.setProp("key2", "val2")

	lines := strings.Split(strings.TrimSuffix(sec.String(), "\n"), "\n")
	if lines[0] != "[sec1]" {
		t.Errorf("header = %q, want [sec1]", lines[0])
	}
	props := lines[1:]
	sort.Strings(props)
	if strings.Join(props, ",") != "key1=val1,key2=val2" {
		t.Errorf("props = %v, want [key1=val1 key2=val2]", props)
	}

	if got := newSection("empty").StringPropsOnly(); got != "" {
		t.Errorf("StringPropsOnly of empty section = %q, want empty", got)
	}
}
