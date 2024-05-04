package newt

import (
	"fmt"
 	"testing"
)

func TestMkName(t *testing.T) {
	mIds := []string{ "people", "groups", "funders", "publishers", "publications" }
	actions := []string { "create", "update", "delete" }
	for _, mId := range mIds {
		for _, action := range actions {
			expected := fmt.Sprintf("%s_%s_form.tmpl", mId, action)
			got := mkName(mId, action, "_form.tmpl")
			if expected != got {
				t.Errorf("expected %q, got %q for %s -> %s", expected, got, mId, action)
			}
		}
	}
}

func TestInList(t *testing.T) {
	l := []string{ "one", "two", "three" }
	for _, item := range l {
		if ! inList(item, l) {
			t.Errorf("expected %q to found in list %+v", item, l)
		}
	}
	item := "fourteen"
	if inList(item, l) {
		t.Errorf("expected to NOT find %q in list %+v", item, l)
	}
}
	
