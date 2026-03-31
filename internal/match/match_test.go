package match

import "testing"

func TestFirstPrefixRule(t *testing.T) {
	rules := []Rule{
		{Match: "https://work.example.com/", Profile: "Work"},
		{Match: "https://", Profile: "Personal"},
	}
	cases := []struct {
		url      string
		wantProf string
		wantOK   bool
	}{
		{"https://work.example.com/foo", "Work", true},
		{"https://other.com/", "Personal", true},
		{"http://local/", "", false},
	}
	for _, tc := range cases {
		p, ok := FirstPrefixRule(tc.url, rules)
		if ok != tc.wantOK || p != tc.wantProf {
			t.Errorf("FirstPrefixRule(%q) = (%q, %v), want (%q, %v)", tc.url, p, ok, tc.wantProf, tc.wantOK)
		}
	}
}
