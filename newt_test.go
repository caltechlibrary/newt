package newt


import (
	"testing"
)


func TestEvalRoute(t *testing.T) {
	routeDefs := []string{
		`/`,
		`/index.html`,
		`/blog/{year year}/{month month}/{day day}`,
		`/index.{ext extname}`,
		`/{p dirname}{filename basename}{ext extname}`,
	}
	paths := map[string][]bool{
		`/`: []bool{ true, false, false, true, false },
		`/index.html`: []bool{ false, true, false, true, true },
		`/index.htm`: []bool{ false, true, false, true, true },
		`/stuff.htm`: []bool{ false, false, false, false, true },
		`/blog/fred/is/it`: []bool{ false, false, false, false, true },
		`/blog/2023/01/02`: []bool{ false, false, true, false, true },
	}
	expectedVars := map[string]map[string]string{
		`/`: map[string]string{},
		`/index.html`: map[string]string{},
		`/index.htm`: map[string]string{},
		`/stuff.htm`: map[string]string{},
		`/blog/fred/is/it`: map[string]string{},
		`/blog/2023/01/02`: map[string]string{
			"year": "2023",
			"month": "01",
			"day": "02",
		},
	}

	for p, expectedOK := range paths {
		for i, defn := range routeDefs {
    		if ok, m := EvalRoute(defn, p); ok == expectedOK[i] {
    			for k, v := range m {
    				if expected, ok := expectedVars[k]; ok {
    					if len(v) == len(expected) {
    						for k2, v2 := range expected {
    							if val, ok := m[k2]; ok {
    								if val != v2 {
    									t.Errorf("expected (rule %d %q) m[%q] -> %q, got value %q", i, defn, k2, v2, val)
    								}
    							} else {
    								t.Errorf("expected (rule %d %q) %q in map, got %+v", i, defn, k2, m)
    							}
    						}
    					} else {
    						t.Errorf("expected (rule %d %q) %+v, got %+v", i, defn, expected, v)
    					}
    				} else {
    					t.Errorf("unexpected (rule %d %q) %q in map %v", i, defn, k, m)
    				}
    			}
    		} else {
    			t.Errorf("expected (rule %d %q) %q to match route %q", i, defn, p, defn)
    		}
		}
	}
}
