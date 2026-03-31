package config

import (
	"reflect"
	"testing"
)

func TestFlattenRules(t *testing.T) {
	rules := []Rule{
		{Matches: []string{"https://a/", "https://b/"}, Profile: "Work"},
		{Match: "https://c/", Profile: "Personal"},
		{Matches: []string{}, Match: "https://solo/", Profile: "Default"},
	}
	got := FlattenRules(rules)
	want := []FlatRule{
		{"https://a/", "Work"},
		{"https://b/", "Work"},
		{"https://c/", "Personal"},
		{"https://solo/", "Default"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FlattenRules() = %#v, want %#v", got, want)
	}
}
